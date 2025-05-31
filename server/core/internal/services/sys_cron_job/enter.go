package sys_cron_job

import (
	"context"
	"encoding/json"
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/pkg/timer"
	"github.com/gin-gonic/gin"
	"time"
)

var controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct {
}

func (_ *Mounter) Initialize() error {
	// 注册一个名为test的定时任务处理函数
	err := global.Timer.RegisterHandler("sys_clear_cron_logs", &timer.HandlerOption{
		Name: "计划任务-删除任务日志",
		Params: timer.HandlerParams{
			{
				Name:     "保留天数",
				Key:      "day",
				Type:     timer.IntParams,
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
			return global.DB.Model(SysCronJobLog{}).Where("created_at < ?", cutoff).Delete(SysCronJobLog{}).Error
		},
	})
	// 初始化数据库
	err = global.DB.AutoMigrate(&SysCronJob{}, &SysCronJobLog{})
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
