package scheduler

import (
	"context"
	"errors"
	"fmt"
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
		options = append(options, gocron.WithLocation(option.Location))
	}
	if option.DistributedLocker != nil {
		options = append(options, gocron.WithDistributedLocker(option.DistributedLocker))
	}
	scheduler, err := gocron.NewScheduler(options...)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &defaultManage{
		ctx:           ctx,
		cancel:        cancel,
		jobSyncer:     option.JobSyncer,
		lock:          sync.Mutex{},
		cronScheduler: scheduler,
		handlers:      make(map[string]*HandlerOption),
		jobs:          make(map[string]Job),
	}, nil
}

type defaultManage struct {
	ctx           context.Context
	cancel        context.CancelFunc
	startOnce     sync.Once
	stopOnce      sync.Once
	wg            sync.WaitGroup
	stopErr       error
	jobSyncer     JobSyncer
	lock          sync.Mutex
	cronScheduler gocron.Scheduler
	handlers      map[string]*HandlerOption
	jobs          map[string]Job
}

// GetHandlers 查询所有已注册处理函数
func (t *defaultManage) GetHandlers() map[string]*HandlerOption {
	t.lock.Lock()
	defer t.lock.Unlock()
	handlers := make(map[string]*HandlerOption, len(t.handlers))
	for key, handler := range t.handlers {
		handlers[key] = handler
	}
	return handlers
}

func (t *defaultManage) getHandler(key string) (*HandlerOption, bool) {
	t.lock.Lock()
	defer t.lock.Unlock()

	handler, ok := t.handlers[key]
	return handler, ok
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
	defer timer.Stop()
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
				delete(t.jobs, jobId)
				// 校验Handler是否存在
				handler, hasHandler := t.getHandler(job.GetHandlerKey())
				if !hasHandler {
					log.Printf("任务 %s 的处理函数 %s 未注册", jobId, job.GetHandlerKey())
					continue
				}
				jobOptions := []gocron.JobOption{
					gocron.WithSingletonMode(gocron.LimitModeReschedule),
					gocron.WithTags(jobId), gocron.WithName(jobId),
				}
				nextRuntime := job.GetNextRunTime()
				if nextRuntime.After(time.Now().UTC()) {
					jobOptions = append(jobOptions, gocron.WithStartAt(gocron.WithStartDateTime(nextRuntime)))
				}
				_, createErr := t.cronScheduler.NewJob(
					gocron.CronJob(job.GetCronExpr(), false),
					gocron.NewTask(func(ctx context.Context) {
						t.execute(ctx, job, handler)
					}),
					jobOptions...,
				)
				if createErr != nil {
					log.Printf("注册任务 %s 失败：%v", jobId, createErr)
					continue
				}
				t.jobs[jobId] = job
			}
			// 删除旧任务
			for oldId, _ := range t.jobs {
				_, has := newJobs[oldId]
				if !has {
					t.cronScheduler.RemoveByTags(oldId)
					delete(t.jobs, oldId)
				}
			}
		case <-t.ctx.Done():
			return
		}
	}
}

// StartScheduler 启动调度器
func (t *defaultManage) StartScheduler(syncInterval time.Duration) {
	t.startOnce.Do(func() {
		t.wg.Add(1)
		go func() { defer t.wg.Done(); t.startSync(syncInterval) }()
		t.cronScheduler.Start()
	})
}

func (t *defaultManage) StopScheduler() error {
	t.stopOnce.Do(func() {
		t.cancel()
		t.wg.Wait()
		t.stopErr = t.cronScheduler.Shutdown()
	})
	return t.stopErr
}

func (t *defaultManage) execute(ctx context.Context, job Job, handler *HandlerOption) {
	run, err := job.Claim(ctx)
	if err != nil {
		log.Printf("领取任务 %s 失败：%v", job.GetID(), err)
		return
	}
	if run == nil {
		return
	}
	var taskErr error
	defer func() {
		if recovered := recover(); recovered != nil {
			taskErr = fmt.Errorf("任务发生 panic：%v", recovered)
		}
		// 停止时任务上下文可能已取消，使用独立超时保存执行结果。
		finishCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := run.Finish(finishCtx, taskErr); err != nil {
			log.Printf("记录任务 %s 执行结果失败：%v", job.GetID(), err)
		}
	}()
	if err := ctx.Err(); err != nil {
		taskErr = err
		return
	}
	taskErr = handler.Handler(ctx, job.GetParams())
}
