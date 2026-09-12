package initialization

import (
	"context"
	"github.com/Gary-Yez/go-admin/internal/cache"
	"github.com/Gary-Yez/go-admin/internal/scheduler"
	"github.com/Gary-Yez/go-admin/internal/system/sys_cron_job"
	"github.com/go-co-op/gocron/v2"
	"time"
)

type cacheLock struct {
	lock cache.Lock
}

func (l *cacheLock) Unlock(_ context.Context) error {
	return l.lock.Unlock()
}

type cacheLocker struct {
	cache cache.Cache
}

func (r *cacheLocker) Lock(ctx context.Context, key string) (gocron.Lock, error) {
	lock, err := r.cache.Lock(ctx, "cron:"+key, time.Second*10, time.Second*5)
	if err != nil {
		return nil, err
	}
	return &cacheLock{lock: lock}, nil
}

func initTaskManager(cacheStore cache.Cache) (scheduler.Scheduler, error) {
	return scheduler.NewManager(scheduler.SchedulerOption{
		JobSyncer:         sys_cron_job.Service,
		DistributedLocker: &cacheLocker{cache: cacheStore},
		Location:          scheduler.DefaultLocation,
	})
}
