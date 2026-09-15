package sys_file

import (
	"context"

	"github.com/Gary-Yez/go-admin/internal/scheduler"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"
	"github.com/gin-gonic/gin"
)

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct{}

func (*Mounter) Name() string { return "文件管理-文件列表" }
func (*Mounter) Initialize() error {
	if err := state.DB().AutoMigrate(&SysFile{}, &SysFilePart{}); err != nil {
		return err
	}
	return state.Scheduler().RegisterHandler("sys_clear_file_uploads", &scheduler.HandlerOption{
		Name: "文件管理-清理过期上传",
		Handler: func(ctx context.Context, _ []byte) error {
			_, err := Service.Cleanup(ctx)
			return err
		},
	})
}
func (*Mounter) Menus() []sys_menu.Definition {
	return []sys_menu.Definition{
		{Name: "文件管理", Key: "sys_files", Icon: "iconoir:folder", Path: "sys_files", Sort: 4},
		{Name: "文件列表", Key: "sys_file", ParentKey: "sys_files", Icon: "iconoir:page", Path: "sys_file", Component: "../core/views/sys_file/index.vue", Sort: 0},
	}
}
func (*Mounter) AdminRouter(group *gin.RouterGroup) {
	group.POST("list", Controller.List)
	group.GET("options", Controller.Options)
	group.POST("upload", Controller.Upload)
	group.POST("begin", Controller.Begin)
	group.GET("session", Controller.Session)
	group.POST("part", Controller.Part)
	group.POST("complete", Controller.Complete)
	group.POST("abort", Controller.Abort)
	group.POST("delete", Controller.Delete)
	group.POST("cleanup", Controller.Cleanup)
	group.POST("link", Controller.Link)
}
func (*Mounter) PublicRouter(group *gin.RouterGroup) {
	// 此地址只接受短期随机下载凭证，不暴露原存储路径和配置。
	group.GET("content/:ticket", Controller.Content)
}
