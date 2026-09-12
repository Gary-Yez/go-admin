package sys_role

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Gary-Yez/go-admin/internal/cache"
	"github.com/Gary-Yez/go-admin/internal/config"
	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/persist"
	rediswatcher "github.com/casbin/redis-watcher/v2"
)

var policyReloadMu sync.Mutex

// StartPolicyRefresh 为遗漏的 Redis 通知提供兜底，返回后台任务的停止函数。
func StartPolicyRefresh() func() {
	return startPolicyRefresh(Enforcer, 30*time.Second)
}

func startPolicyRefresh(enforcer *casbin.SyncedCachedEnforcer, interval time.Duration) func() {
	stop := make(chan struct{})
	done := make(chan struct{})
	var once sync.Once
	go func() {
		defer close(done)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				if err := reloadPolicy(enforcer); err != nil {
					log.Printf("定期同步角色权限失败，将在下个周期重试：%v", err)
				}
			}
		}
	}()
	return func() {
		once.Do(func() { close(stop) })
		<-done
	}
}

func policyChannel(cfg *config.Config) string {
	return cache.DatabaseNamespace(cfg.Mysql.Host, cfg.Mysql.Port, cfg.Mysql.Database) + "channel:casbin"
}

func newPolicyWatcher(enforcer *casbin.SyncedCachedEnforcer, cfg *config.Config) (persist.Watcher, error) {
	watcher, err := rediswatcher.NewWatcher(cfg.Redis.Address(), rediswatcher.WatcherOptions{
		Options:    *cfg.Redis.Option(),
		Channel:    policyChannel(cfg),
		IgnoreSelf: true,
		OptionalUpdateCallback: func(string) {
			// 从已提交的数据库重新加载，避免增量通知留下旧的鉴权缓存。
			if err := reloadPolicy(enforcer); err != nil {
				log.Printf("同步角色权限失败：%v", err)
			}
		},
	})
	if err != nil {
		return nil, err
	}
	if err := enforcer.SetWatcher(watcher); err != nil {
		watcher.Close()
		return nil, err
	}
	return watcher, nil
}

func reloadPolicy(enforcer *casbin.SyncedCachedEnforcer) error {
	policyReloadMu.Lock()
	defer policyReloadMu.Unlock()
	err := enforcer.LoadPolicy()
	return errors.Join(err, enforcer.InvalidateCache())
}

// 数据库事务提交后重新加载权限并通知其他实例，不再重复写入策略。
func (s *serviceStruck) ReloadPolicy() error {
	err := reloadPolicy(Enforcer)
	if policyWatcher != nil {
		// 即使本机加载失败，也通知其他实例刷新已提交的数据。
		err = errors.Join(err, policyWatcher.Update())
	}
	return err
}
