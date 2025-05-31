package core

import (
	"cmp"
	"fmt"
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
	"gitee.com/mxcker/go-admin/server/pkg/modular"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"slices"
	"strings"
	"time"
)

func Run(needInit bool) {
	// 注册系统组件
	loader := modular.NewLoader()
	// 只有开发环境下注册自动生成代码
	if global.Config.IsDev() {
		loader.Mount("sys_devtools", new(sys_devtools.Mounter))
	}
	loader.Mount("sys_menu", new(sys_menu.Mounter))
	loader.Mount("sys_apis", new(sys_apis.Mounter))
	loader.Mount("sys_role", new(sys_role.Mounter))
	loader.Mount("sys_admin", new(sys_admin.Mounter))
	loader.Mount("sys_auth", new(sys_auth.Mounter))
	loader.Mount("sys_cron_job", new(sys_cron_job.Mounter))
	//加载用户组件
	err := modules.Load(loader)
	if err != nil {
		panic(err)
	}
	if needInit {
		err = loader.InitializeAll()
		if err != nil {
			panic(err)
		}
	}
	//创建服务器
	server := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	server.Use(cors.New(corsConfig))
	server.Use(static.Serve(global.Config.Server.AdminPrefix, static.LocalFile("./dist", false)))
	var apiPrefix = strings.TrimSpace(global.Config.Server.ApiPrefix)
	AdminAuthGroup := server.Group(apiPrefix, middlewares.JWTMiddleware(), middlewares.CasbinMiddleware(sys_role.Enforcer))
	PublicGroup := server.Group(apiPrefix)
	err = loader.RegisterRouter(AdminAuthGroup, PublicGroup)
	if err != nil {
		panic(err)
	}
	global.Routes = server.Routes()
	slices.SortFunc(global.Routes, func(a, b gin.RouteInfo) int {
		return cmp.Compare(a.Path, b.Path)
	})
	// 在系统成功挂载后再启动调度器
	global.Timer.StartScheduler(time.Second)
	// 监听并启动服务
	listenAddr := viper.GetString("server.host") + ":" + viper.GetString("server.port")
	fmt.Println("程序运行在：" + listenAddr)
	fmt.Println("管理员入口：" + global.Config.Server.AdminPrefix)
	err = server.Run(listenAddr)
	if err != nil {
		panic(err)
	}
}
