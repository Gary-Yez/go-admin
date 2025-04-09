package core

import (
	"gitee.com/mxcker/go-admin/server/core/global"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_admin"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_auth"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_autocode"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_global_variable"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_home"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_menu"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_role"
	"gitee.com/mxcker/go-admin/server/core/middlewares"
	"gitee.com/mxcker/go-admin/server/core/models"
	"gitee.com/mxcker/go-admin/server/modules"
	"github.com/gin-gonic/gin"
)

func Load(Server *gin.Engine, needInit bool) {
	AdminAuthGroup := Server.Group("/api", middlewares.SysAuth)
	PublicGroup := Server.Group("/api")
	// 注册系统组件
	mounterMap := new(models.MounterMap)
	// 只有开发环境下注册自动生成代码
	if global.IsDevelopment {
		mounterMap.Add("sys/autocode", new(sys_autocode.Mounter))
	}
	mounterMap.Add("sys/global_variable", new(sys_global_variable.Mounter))
	mounterMap.Add("sys/menu", new(sys_menu.Mounter))
	mounterMap.Add("sys/role", new(sys_role.Mounter))
	mounterMap.Add("sys/admin", new(sys_admin.Mounter))
	mounterMap.Add("sys/auth", new(sys_auth.Mounter))
	mounterMap.Add("sys/home", new(sys_home.Mounter))
	//加载用户组件
	err := modules.Load(mounterMap)
	if err != nil {
		panic(err)
	}
	if needInit {
		err = mounterMap.InitAll()
		if err != nil {
			panic(err)
		}
	}
	err = mounterMap.RegisterAll(AdminAuthGroup, PublicGroup)
	if err != nil {
		panic(err)
	}
}
