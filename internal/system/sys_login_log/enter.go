package sys_login_log

import (
	"context"
	"encoding/json"
	"github.com/Gary-Yez/go-admin/internal/scheduler"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"
	"github.com/gin-gonic/gin"
)

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct{}

func (*Mounter) Name() string { return "系统运维-登录日志" }
func (*Mounter) Initialize() error {
	if err := state.DB().AutoMigrate(&SysLoginLog{}); err != nil {
		return err
	}
	return state.Scheduler().RegisterHandler("sys_clear_login_logs", &scheduler.HandlerOption{
		Name: "登录日志-清理历史记录",
		Params: scheduler.HandlerParams{
			{Name: "保留天数", Key: "days", Type: scheduler.IntParams, Required: true, Description: "保留最近多少天的登录日志，范围 1–36500"},
		},
		Handler: func(ctx context.Context, params []byte) error {
			body := new(CleanupBody)
			if err := json.Unmarshal(params, body); err != nil {
				return err
			}
			_, err := Service.Cleanup(ctx, body.Days)
			return err
		},
	})
}
func (*Mounter) Menus() []sys_menu.Definition {
	return []sys_menu.Definition{{Name: "登录日志", Key: "sys_login_log", ParentKey: "sys_operations", Icon: "iconoir:log-in", Path: "sys_login_log", Component: "../core/views/sys_login_log/index.vue", Sort: 1}}
}
func (*Mounter) AdminRouter(group *gin.RouterGroup) {
	group.POST("list", Controller.List)
	group.POST("delete", Controller.Delete)
	group.POST("cleanup", Controller.Cleanup)
}
func (*Mounter) PublicRouter(group *gin.RouterGroup) {}
