package core

import (
	"gitee.com/mxcker/go-admin/server/core/global"
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_admin"
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_auth"
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_autocode"
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_cron_job"
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_menu"
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_role"
	"gitee.com/mxcker/go-admin/server/core/middlewares"
	"gitee.com/mxcker/go-admin/server/core/pkg/modular"
	"gitee.com/mxcker/go-admin/server/modules"
	"github.com/gin-gonic/gin"
	"time"
)

func Load(Server *gin.Engine, needInit bool) {
	AdminAuthGroup := Server.Group("/api", middlewares.SysAuth)
	PublicGroup := Server.Group("/api")
	// 注册系统组件
	loader := modular.NewLoader()
	// 只有开发环境下注册自动生成代码
	if global.IsDev() {
		loader.Add("sys/autocode", new(sys_autocode.Mounter))
	}
	loader.Add("sys/menu", new(sys_menu.Mounter))
	loader.Add("sys/role", new(sys_role.Mounter))
	loader.Add("sys/admin", new(sys_admin.Mounter))
	loader.Add("sys/auth", new(sys_auth.Mounter))
	loader.Add("sys/cron_job", new(sys_cron_job.Mounter))
	//加载用户组件
	err := modules.Load(loader)
	if err != nil {
		panic(err)
	}
	if needInit {
		err = loader.Initialize()
		if err != nil {
			panic(err)
		}
	}
	err = loader.Server(AdminAuthGroup, PublicGroup)
	if err != nil {
		panic(err)
	}
	// 在系统成功挂载后再启动调度器
	global.Timer.StartScheduler(time.Second)
}
