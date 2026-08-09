package sys_cron_job

import (
	"context"
	"encoding/json"
	"github.com/Gary-Yez/go-admin/internal/state"

	"github.com/Gary-Yez/go-admin/scheduler"
	"github.com/gin-gonic/gin"
	"time"
)

var controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct {
}

func (_ *Mounter) Name() string {
	return "核心服务-定时任务"
}

func (_ *Mounter) Initialize() error {
	// 注册一个名为test的定时任务处理函数
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
			cutoff := time.Now().AddDate(0, 0, -data.Day)
			return state.DB().Model(SysCronJobLog{}).Where("created_at < ?", cutoff).Delete(SysCronJobLog{}).Error
		},
	})
	// 初始化数据库
	err = state.DB().AutoMigrate(&SysCronJob{}, &SysCronJobLog{})
	if err != nil {
		return err
	}
	return nil
}

func (_ *Mounter) AdminRouter(adminAuthGroup *gin.RouterGroup) {
	adminAuthGroup.GET("get_handlers", controller.GetHandlers)
	adminAuthGroup.GET("get", controller.Get)
	adminAuthGroup.POST("list", controller.List)
	adminAuthGroup.POST("logs", controller.GetLogs)
	adminAuthGroup.POST("create", controller.Create)
	adminAuthGroup.POST("delete", controller.Delete)
	adminAuthGroup.POST("edit", controller.Edit)
}

func (_ *Mounter) PublicRouter(publicGroup *gin.RouterGroup) {

}
