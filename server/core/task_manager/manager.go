package task_manager

import (
	"context"
	"errors"
	"github.com/go-co-op/gocron/v2"
	"strconv"
	"sync"
)

func NewManager() Manager {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}
	return &defaultManage{
		lock:          sync.Mutex{},
		cronScheduler: scheduler,
		handlers:      make(map[string]*HandlerOption),
		jobs:          make(map[uint]*ScheduledJob),
	}
}

type defaultManage struct {
	lock          sync.Mutex
	cronScheduler gocron.Scheduler
	handlers      map[string]*HandlerOption
	jobs          map[uint]*ScheduledJob
}

// GetHandlers 查询所有已注册处理函数
func (t *defaultManage) GetHandlers() map[string]*HandlerOption {
	return t.handlers
}

// RegisterHandler 注册处理函数
func (t *defaultManage) RegisterHandler(
	key string,
	option *HandlerOption,
) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	_, ok := t.handlers[key]
	if ok {
		return errors.New("TaskKey already registered")
	}
	t.handlers[key] = option
	return nil
}

// RegisterJob 注册定时任务
func (t *defaultManage) RegisterJob(job *ScheduledJob) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	// 校验Handler是否存在
	handler, ok := t.handlers[job.HandlerKey]
	if !ok {
		return errors.New("handler " + job.HandlerKey + " not registered")
	}
	// 校验定时任务是否已存在
	if _, exists := t.jobs[job.Id]; exists {
		return errors.New("schedule Id already exists")
	}
	if job.ctx == nil {
		ctx, cancel := context.WithCancel(context.Background())
		job.ctx = ctx
		job.cancel = cancel
	}
	var jobOptions []gocron.JobOption
	if !job.AllowConcurrency {
		jobOptions = append(jobOptions, gocron.WithSingletonMode(gocron.LimitModeReschedule))
	}
	jobOptions = append(jobOptions, gocron.WithContext(job.ctx))
	jobOptions = append(jobOptions, gocron.WithTags(strconv.Itoa(int(job.Id))))
	// 添加任务到cron
	_, err := t.cronScheduler.NewJob(gocron.CronJob(job.CronExpr, true), gocron.NewTask(func(ctx context.Context) {
		// 执行任务处理函数
		_ = handler.Handler(ctx, job.Params)
	}), jobOptions...)
	if err != nil {
		return err
	}
	t.jobs[job.Id] = job
	return nil
}

// RemoveJobs 删除定时任务
func (t *defaultManage) RemoveJobs(ids []uint) {
	t.lock.Lock()
	defer t.lock.Unlock()
	var tags []string
	for _, id := range ids {
		if _, ok := t.jobs[id]; ok {
			t.jobs[id].cancel()
		}
		tags = append(tags, strconv.Itoa(int(id)))
		delete(t.jobs, id)
	}
	t.cronScheduler.RemoveByTags(tags...)
	return
}

// StartScheduler 启动调度器
func (t *defaultManage) StartScheduler() {
	t.cronScheduler.Start()
}

// StopScheduler 停止调度器
func (t *defaultManage) StopScheduler() {
	//t.cronScheduler.Stop()
}
