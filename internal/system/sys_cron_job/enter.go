package sys_cron_job

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"

	"github.com/Gary-Yez/go-admin/internal/scheduler"
	"github.com/gin-gonic/gin"
	"time"
)

var controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct {
}

func (_ *Mounter) Name() string {
	return "系统运维-计划任务"
}

func (_ *Mounter) Menus() []sys_menu.Definition {
	return []sys_menu.Definition{
		{
			Name: "计划任务", Key: "sys_cron_job", ParentKey: "sys_operations", Icon: "iconoir:clock",
			Path: "sys_cron_job", Component: "../core/views/sys_cron_job/index.vue", Sort: 2,
		},
	}
}

func (_ *Mounter) Initialize() error {
	// 注册任务日志清理处理函数。
	err := state.Scheduler().RegisterHandler("sys_clear_cron_logs", &scheduler.HandlerOption{
		Name: "计划任务-删除任务日志",
		Params: scheduler.HandlerParams{
			{
				Name:     "保留天数",
				Key:      "day",
				Type:     scheduler.IntParams,
				Required: true,
			},
		},
		Handler: func(ctx context.Context, params []byte) error {
			data := new(struct {
				Day int `json:"day"`
			})
			err := json.Unmarshal(params, &data)
			if err != nil {
				return err
			}
			if data.Day <= 0 {
				return fmt.Errorf("保留天数必须大于 0")
			}
			cutoff := time.Now().AddDate(0, 0, -data.Day)
			return state.DB().WithContext(ctx).Where("created_at < ? AND status IN ?", cutoff, []string{"success", "failed", "interrupted"}).Delete(&SysCronJobLog{}).Error
		},
	})
	if err != nil {
		return err
	}
	// 初始化数据库
	err = state.DB().AutoMigrate(&SysCronJob{}, &SysCronJobLog{}, &initializationRecord{})
	if err != nil {
		return err
	}
	return InitData()
}

func (_ *Mounter) AdminRouter(adminAuthGroup *gin.RouterGroup) {
	adminAuthGroup.GET("get_handlers", controller.GetHandlers)
	adminAuthGroup.POST("list", controller.List)
	adminAuthGroup.POST("logs", controller.GetLogs)
	adminAuthGroup.POST("create", controller.Create)
	adminAuthGroup.POST("delete", controller.Delete)
	adminAuthGroup.POST("edit", controller.Edit)
}

func (_ *Mounter) PublicRouter(publicGroup *gin.RouterGroup) {

}
