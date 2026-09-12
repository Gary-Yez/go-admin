package sys_monitor

import (
	"context"
	"errors"
	"github.com/Gary-Yez/go-admin/internal/monitor"
)

type serviceStruct struct{}

var store *monitor.Store

// SetStore 仅在启动阶段设置，采集循环与接口共用同一存储。
func SetStore(value *monitor.Store) { store = value }
func (*serviceStruct) List(ctx context.Context) ([]monitor.Snapshot, error) {
	if store == nil {
		return nil, errors.New("节点监控尚未初始化")
	}
	return store.List(ctx)
}
