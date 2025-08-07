package internal

import (
	"cmp"
	"errors"
	"gitee.com/mxcker/go-admin/server/core"
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_admin"
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_apis"
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_auth"
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_cron_job"
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_devtools"
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_menu"
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_role"
	"gitee.com/mxcker/go-admin/server/core/middlewares"
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/modules"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"slices"
	"strings"
	"time"
)

func Run(needInit bool) {
	// 注册系统组件
	if global.Config.IsDev() {
		core.AddModule("sys_devtools", new(sys_devtools.Mounter))
	}
	core.AddModule("sys_menu", new(sys_menu.Mounter))
	core.AddModule("sys_apis", new(sys_apis.Mounter))
	core.AddModule("sys_role", new(sys_role.Mounter))
	core.AddModule("sys_admin", new(sys_admin.Mounter))
	core.AddModule("sys_auth", new(sys_auth.Mounter))
	core.AddModule("sys_cron_job", new(sys_cron_job.Mounter))
	//加载用户组件
	modules.Init()
	// 初始化模块
	for _, module := range core.Modules {
		log.Println("初始化模块：", module.Name())
		if err := module.Initialize(); err != nil {
			panic(errors.New("initialize fail: " + err.Error()))
		}
	}
	//创建服务器
	if global.Config.IsDev() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	server := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	server.Use(cors.New(corsConfig))
	server.Use(static.Serve(global.Config.Server.AdminPrefix, static.LocalFile("./dist", false)))
	var apiPrefix = strings.TrimSpace(global.Config.Server.ApiPrefix)
	AdminAuthGroup := server.Group(apiPrefix, middlewares.JWTMiddleware(), middlewares.CasbinMiddleware(sys_role.Enforcer))
	PublicGroup := server.Group(apiPrefix)
	// 注册路由
	for routePath, module := range core.Modules {
		log.Println("注册路由中：", module.Name())
		module.AdminRouter(AdminAuthGroup.Group(routePath))
		module.PublicRouter(PublicGroup.Group(routePath))
	}
	global.Routes = server.Routes()
	slices.SortFunc(global.Routes, func(a, b gin.RouteInfo) int {
		return cmp.Compare(a.Path, b.Path)
	})
	// 在系统成功挂载后再启动调度器
	global.Timer.StartScheduler(time.Second)
	// 监听并启动服务
	listenAddr := viper.GetString("server.host") + ":" + viper.GetString("server.port")
	log.Println("程序运行在：" + listenAddr)
	log.Println("管理员入口：" + global.Config.Server.AdminPrefix)
	err := server.Run(listenAddr)
	if err != nil {
		panic(err)
	}
}
