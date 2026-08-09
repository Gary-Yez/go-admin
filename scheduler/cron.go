package scheduler

import (
	"context"
	"errors"
	"github.com/go-co-op/gocron/v2"
	"log"
	"sync"
	"time"
)

func NewManager(option SchedulerOption) (Scheduler, error) {
	if option.JobSyncer == nil {
		return nil, errors.New("JobSyncer is nil")
	}
	var options []gocron.SchedulerOption
	if option.Location != nil {
		options = append(options, gocron.WithLocation(time.UTC))
	}
	if option.DistributedLocker != nil {
		options = append(options, gocron.WithDistributedLocker(option.DistributedLocker))
	}
	scheduler, err := gocron.NewScheduler(options...)
	if err != nil {
		panic(err)
	}
	return &defaultManage{
		ctx:           context.Background(),
		jobSyncer:     option.JobSyncer,
		lock:          sync.Mutex{},
		cronScheduler: scheduler,
		handlers:      make(map[string]*HandlerOption),
		jobs:          make(map[string]Job),
	}, nil
}

type defaultManage struct {
	ctx           context.Context
	jobSyncer     JobSyncer
	lock          sync.Mutex
	cronScheduler gocron.Scheduler
	handlers      map[string]*HandlerOption
	jobs          map[string]Job
}

// GetHandlers 查询所有已注册处理函数
func (t *defaultManage) GetHandlers() map[string]*HandlerOption {
	return t.handlers
}

// RegisterHandler 注册处理函数
func (t *defaultManage) RegisterHandler(key string, option *HandlerOption) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	_, ok := t.handlers[key]
	if ok {
		return errors.New("处理函数Key已被注册")
	}
	t.handlers[key] = option
	return nil
}

func (t *defaultManage) startSync(syncInterval time.Duration) {
	timer := time.NewTicker(syncInterval)
	for {
		select {
		case <-timer.C:
			jobs, err := t.jobSyncer.Sync()
			if err != nil {
				log.Println("定时任务调度器同步任务失败:", err)
				continue
			}
			newJobs := make(map[string]Job, len(jobs))
			//循环job进行处理
			for _, job := range jobs {
				jobId := job.GetID()
				newJobs[jobId] = job
				// 如果旧的任务是否存在或版本未发生改变则跳出循环
				if oldJob, ok := t.jobs[jobId]; ok && oldJob.GetVersion() == job.GetVersion() {
					continue
				}
				t.cronScheduler.RemoveByTags(jobId)
				// 校验Handler是否存在
				handler, hasHandler := t.handlers[job.GetHandlerKey()]
				if !hasHandler {
					continue
				}
				var startTimeOption gocron.JobOption
				nextRuntime := job.GetNextRunTime()
				if nextRuntime.Before(time.Now().UTC()) {
					startTimeOption = gocron.WithStartAt(gocron.WithStartImmediately())
				} else {
					startTimeOption = gocron.WithStartAt(gocron.WithStartDateTime(nextRuntime))
				}
				_, createErr := t.cronScheduler.NewJob(
					gocron.CronJob(job.GetCronExpr(), false),
					gocron.NewTask(func(ctx context.Context) {
						startTime := time.Now().UTC()
						// 执行任务处理函数
						taskErr := handler.Handler(ctx, job.GetParams())
						if !errors.Is(taskErr, context.Canceled) {
							endTime := time.Now().UTC()
							// 更新下次运行时间
							var nextTime = job.GetNextRunTime()
							for _, j := range t.cronScheduler.Jobs() {
								if j.Name() == job.GetID() {
									if next, gerErr := j.NextRun(); gerErr == nil {
										nextTime = next
									}
									break
								}
							}
							job.AfterRun(startTime, endTime, nextTime, taskErr)
						}
					}),
					gocron.WithSingletonMode(gocron.LimitModeReschedule),
					gocron.WithTags(jobId),
					gocron.WithName(jobId),
					startTimeOption,
				)
				if createErr != nil {
					continue
				}
				t.jobs[jobId] = job
			}
			// 删除旧任务
			for oldId, _ := range t.jobs {
				_, has := newJobs[oldId]
				if !has {
					t.cronScheduler.RemoveByTags(oldId)
				}
			}
		case <-t.ctx.Done():
			return
		}
	}
}

// StartScheduler 启动调度器
func (t *defaultManage) StartScheduler(syncInterval time.Duration) {
	go t.startSync(syncInterval)
	t.cronScheduler.Start()
}
