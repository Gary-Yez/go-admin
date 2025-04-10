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
	"gitee.com/mxcker/go-admin/server/core/internal/sys_task"
	"gitee.com/mxcker/go-admin/server/core/middlewares"
	"gitee.com/mxcker/go-admin/server/core/types/module_loader"
	"gitee.com/mxcker/go-admin/server/modules"
	"github.com/gin-gonic/gin"
)

func Load(Server *gin.Engine, needInit bool) {
	AdminAuthGroup := Server.Group("/api", middlewares.SysAuth)
	PublicGroup := Server.Group("/api")
	// 注册系统组件
	loader := module_loader.NewLoader()
	// 只有开发环境下注册自动生成代码
	if global.IsDevelopment {
		loader.Add("sys/autocode", new(sys_autocode.Mounter))
	}
	loader.Add("sys/global_variable", new(sys_global_variable.Mounter))
	loader.Add("sys/menu", new(sys_menu.Mounter))
	loader.Add("sys/role", new(sys_role.Mounter))
	loader.Add("sys/admin", new(sys_admin.Mounter))
	loader.Add("sys/auth", new(sys_auth.Mounter))
	loader.Add("sys/home", new(sys_home.Mounter))
	loader.Add("sys/task", new(sys_task.Mounter))
	//加载用户组件
	err := modules.Load(loader)
	if err != nil {
		panic(err)
	}
	if needInit {
		err = loader.InitAll()
		if err != nil {
			panic(err)
		}
	}
	err = loader.RegisterAll(AdminAuthGroup, PublicGroup)
	if err != nil {
		panic(err)
	}
	// 在系统成功挂载后再启动调度器
	global.TaskManager.StartScheduler()
}
