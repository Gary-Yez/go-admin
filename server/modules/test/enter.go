// Package test 这是一个演示模块
package test

import (
	"context"
	"fmt"
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
	err := global.Timer.RegisterHandler("test", &timer.HandlerOption{
		Name: "测试任务",
		Params: timer.HandlerParams{
			{
				Name:     "a",
				Key:      "a",
				Type:     "string",
				Required: true,
			},
		},
		Handler: func(ctx context.Context, params []byte) error {
			fmt.Println("test任务开始了:", string(params))
			ticker := time.NewTicker(time.Second * 5)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-ticker.C:
					return nil
				}
			}
		},
	})
	if err != nil {
		return err
	}
	// 这里执行一些初始化操作
	return nil
}

func (_ *Mounter) Register(path string, adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error {
	// 注册路由
	// 需要管理员权限的路由
	Group := adminAuthGroup.Group(path)
	{
		_ = Group
	}
	// 无需鉴权的路由
	Public := publicGroup.Group(path)
	{
		_ = Public
	}
	return nil
}
