package sys_api_token

import (
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"
	"github.com/Gary-Yez/go-admin/response"
	"github.com/gin-gonic/gin"
)

var Controller = new(controllerStruct)
var Service = new(serviceStruct)

type Mounter struct{}

func (*Mounter) Name() string { return "权限管理-API密钥管理" }

func (*Mounter) Initialize() error {
	if err := state.DB().AutoMigrate(&SysApiToken{}); err != nil {
		return err
	}
	// 页面归入同一模块，仅更新原组件路径，保留菜单配置和角色关联。
	return state.DB().Model(&sys_menu.SysMenu{}).
		Where(map[string]any{"key": "sys_api_token", "component": "../core/views/sys_api_token_manage/index.vue"}).
		Update("component", "../core/views/sys_api_token/manage.vue").Error
}
func (*Mounter) Menus() []sys_menu.Definition {
	return []sys_menu.Definition{{Name: "API密钥管理", Key: "sys_api_token", ParentKey: "sys_permission", Icon: "iconoir:key", Path: "sys_api_token_manage", Component: "../core/views/sys_api_token/manage.vue", Sort: 4}}
}
func (*Mounter) AdminRouter(group *gin.RouterGroup) {
	group.POST("list", Controller.List)
	group.POST("delete", Controller.Delete)
	mine := group.Group("mine")
	mine.Use(func(ctx *gin.Context) {
		if ctx.GetBool("api_token_auth") {
			response.Error(ctx, "请使用登录会话管理 API 密钥", 403)
		}
	})
	mine.GET("list", Controller.ListApiTokens)
	mine.POST("save", Controller.SaveApiToken)
	mine.POST("delete", Controller.DeleteApiToken)
}
func (*Mounter) PublicRouter(group *gin.RouterGroup) {}
