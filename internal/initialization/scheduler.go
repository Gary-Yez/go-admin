package initialization

import (
	"context"
	"github.com/Gary-Yez/go-admin/cache"
	"github.com/Gary-Yez/go-admin/internal/system/sys_cron_job"
	"github.com/Gary-Yez/go-admin/scheduler"
	"github.com/go-co-op/gocron/v2"
	"time"
)

type cacheLock struct {
	lock cache.Lock
}

func (l *cacheLock) Unlock(_ context.Context) error {
	err := l.lock.Unlock()
	return err
}

type cacheLocker struct {
	cache cache.Cache
}

func (r *cacheLocker) Lock(_ context.Context, key string) (gocron.Lock, error) {
	lock, err := r.cache.Lock("core:cron:"+key, time.Second*10, time.Second*5)
	if err != nil {
		return nil, err
	}
	l := &cacheLock{
		lock: lock,
	}
	return l, nil
}

func initTaskManager(cacheStore cache.Cache) (scheduler.Scheduler, error) {
	t, err := scheduler.NewManager(scheduler.SchedulerOption{
		JobSyncer:         sys_cron_job.Service,
		DistributedLocker: &cacheLocker{cache: cacheStore},
		Location:          time.UTC,
	})
	if err != nil {
		return nil, err
	}
	return t, nil
}
