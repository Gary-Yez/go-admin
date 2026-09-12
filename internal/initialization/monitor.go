package initialization

import (
	"context"
	"github.com/Gary-Yez/go-admin/internal/cache"
	"github.com/Gary-Yez/go-admin/internal/monitor"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_monitor"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
	"time"
)

// StartMonitor 每个实例独立采集，无需抢锁；退出时取消并等待采集循环结束。
func StartMonitor() func() {
	cfg := state.Config()
	var client *redis.Client
	if cfg.Redis.IsNotEmpty() {
		client = state.Cache().Client().(*redis.Client)
	}
	store := monitor.NewStore(client, cache.DatabaseNamespace(cfg.Mysql.Host, cfg.Mysql.Port, cfg.Mysql.Database))
	sys_monitor.SetStore(store)
	collector := monitor.NewCollector(cfg.Server.NodeName)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	var once sync.Once
	go func() {
		defer close(done)
		ticker := time.NewTicker(monitor.Interval)
		defer ticker.Stop()
		failed := false
		for {
			sampleCtx, stop := context.WithTimeout(ctx, 4*time.Second)
			row := collector.Collect(sampleCtx)
			err := store.Publish(sampleCtx, row)
			stop()
			if err != nil && ctx.Err() == nil {
				if !failed {
					log.Printf("节点监控上报失败，将自动重试：%v", err)
				}
				failed = true
			} else if failed && ctx.Err() == nil {
				log.Print("节点监控上报已恢复")
				failed = false
			}
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()
	return func() { once.Do(cancel); <-done }
}
