package sys_cron_job

import (
	"context"
	"encoding/json"
	"gitee.com/mxcker/go-admin/server/core/global"
	"gitee.com/mxcker/go-admin/server/core/pkg/timer"
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

func (_ *Mounter) Register(path string, adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error {
	// 注册路由
	// 需要管理员权限的路由
	Group := adminAuthGroup.Group(path)
	{
		//Group.GET("get_registered_task", controller.GetRegisteredTask)
		Group.GET("get", controller.Get)
		Group.POST("list", controller.List)
		Group.POST("logs", controller.GetLogs)
		Group.POST("create", controller.Create)
		Group.POST("delete", controller.Delete)
		Group.POST("edit", controller.Edit)
	}
	// 无需鉴权的路由
	Public := publicGroup.Group(path)
	{
		_ = Public
		Public.GET("get_registered_handler", controller.GetRegisteredHandler)

	}
	return nil
}
