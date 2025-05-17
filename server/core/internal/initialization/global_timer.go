package initialization

import (
	"context"
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_cron_job"
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/pkg/cacher"
	timer2 "gitee.com/mxcker/go-admin/server/pkg/timer"
	"github.com/go-co-op/gocron/v2"
	"time"
)

type cacheLock struct {
	lock cacher.Lock
}

func (l *cacheLock) Unlock(_ context.Context) error {
	err := l.lock.Unlock()
	return err
}

type cacheLocker struct {
}

func (r *cacheLocker) Lock(_ context.Context, key string) (gocron.Lock, error) {
	lock, err := global.Cache.Lock("core:cron:"+key, time.Second*10, time.Second*5)
	if err != nil {
		return nil, err
	}
	l := &cacheLock{
		lock: lock,
	}
	return l, nil
}

func initTaskManager() error {
	t, err := timer2.NewManager(timer2.SchedulerOption{
		JobSyncer:         sys_cron_job.Service,
		DistributedLocker: &cacheLocker{},
		Location:          time.UTC,
	})
	if err != nil {
		return err
	}
	global.Timer = t
	return nil
}
