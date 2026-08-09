// Package admin provides the complete go-admin runtime.
package admin

import (
	"cmp"
	"errors"
	"fmt"
	"github.com/Gary-Yez/go-admin/internal/initialization"
	"github.com/Gary-Yez/go-admin/internal/middlewares"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_admin"
	"github.com/Gary-Yez/go-admin/internal/system/sys_apis"
	"github.com/Gary-Yez/go-admin/internal/system/sys_auth"
	"github.com/Gary-Yez/go-admin/internal/system/sys_cron_job"
	"github.com/Gary-Yez/go-admin/internal/system/sys_devtools"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"
	"github.com/Gary-Yez/go-admin/internal/system/sys_role"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"log"
	"os"
	"slices"
	"strings"
	"time"
)

type registeredModule struct {
	key string
	Module
}

var registry = struct {
	modules     []registeredModule
	keys        map[string]struct{}
	initialized bool
}{keys: make(map[string]struct{})}

// Register adds an application module to the default admin runtime.
func Register(key string, module Module) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return errors.New("module key cannot be empty")
	}
	if module == nil {
		return fmt.Errorf("module %q cannot be nil", key)
	}
	if registry.initialized {
		return fmt.Errorf("module %q cannot be registered after admin startup", key)
	}
	if _, exists := registry.keys[key]; exists {
		return fmt.Errorf("module %q is already registered", key)
	}
	registry.keys[key] = struct{}{}
	registry.modules = append(registry.modules, registeredModule{key: key, Module: module})
	return nil
}

func MustRegister(key string, module Module) {
	if err := Register(key, module); err != nil {
		panic(err)
	}
}

// Run initializes the singleton admin runtime and starts its HTTP server.
// Development mode loads config.dev.yaml by default; production mode loads config.yaml.
func Run(dev bool) error {
	configFile := "config.yaml"
	if dev {
		configFile = "config.dev.yaml"
	}
	flags := pflag.NewFlagSet("go-admin", pflag.ContinueOnError)
	flags.String("server.host", "0.0.0.0", "Web服务运行的IP")
	flags.String("server.port", "8080", "Web服务运行的端口")
	if err := flags.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("parse arguments: %w", err)
	}
	cfg, err := initialization.InitConfig(flags, configFile)
	if err != nil {
		return err
	}
	cfg.Server.Dev = dev

	deps, err := initialization.InitDependencies(cfg)
	if err != nil {
		return err
	}
	state.Configure(cfg, deps.DB, deps.Cache, deps.Scheduler)

	if err := registerBuiltins(); err != nil {
		return err
	}
	for _, module := range registry.modules {
		log.Println("初始化模块：", module.Name())
		if err := module.Initialize(); err != nil {
			return fmt.Errorf("initialize module %q: %w", module.key, err)
		}
	}
	registry.initialized = true

	if cfg.IsDev() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	engine.Use(cors.New(corsConfig))
	engine.Use(static.Serve(cfg.Server.AdminPrefix, static.LocalFile("./dist", false)))

	apiPrefix := strings.TrimSpace(cfg.Server.ApiPrefix)
	adminGroup := engine.Group(apiPrefix,
		middlewares.JWTMiddleware(cfg.Jwt.Secret, sys_auth.Service.VerifyApiToken),
		middlewares.CasbinMiddleware(apiPrefix, sys_role.Enforcer),
	)
	publicGroup := engine.Group(apiPrefix)
	for _, module := range registry.modules {
		module.AdminRouter(adminGroup.Group(module.key))
		module.PublicRouter(publicGroup.Group(module.key))
	}
	routes := engine.Routes()
	slices.SortFunc(routes, func(left, right gin.RouteInfo) int {
		return cmp.Compare(left.Path, right.Path)
	})
	state.SetRoutes(routes)
	state.Scheduler().StartScheduler(time.Second)

	address := cfg.Server.Host + ":" + cfg.Server.Port
	log.Println("程序运行在：" + address)
	log.Println("管理员入口：" + cfg.Server.AdminPrefix)
	return engine.Run(address)
}

func registerBuiltins() error {
	applicationModuleCount := len(registry.modules)
	if state.Config().IsDev() {
		if err := Register("sys_devtools", new(sys_devtools.Mounter)); err != nil {
			return err
		}
	}
	modules := []registeredModule{
		{key: "sys_menu", Module: new(sys_menu.Mounter)},
		{key: "sys_apis", Module: new(sys_apis.Mounter)},
		{key: "sys_role", Module: new(sys_role.Mounter)},
		{key: "sys_admin", Module: new(sys_admin.Mounter)},
		{key: "sys_auth", Module: new(sys_auth.Mounter)},
		{key: "sys_cron_job", Module: new(sys_cron_job.Mounter)},
	}
	for _, module := range modules {
		if err := Register(module.key, module.Module); err != nil {
			return err
		}
	}
	applicationModules := append([]registeredModule(nil), registry.modules[:applicationModuleCount]...)
	builtinModules := append([]registeredModule(nil), registry.modules[applicationModuleCount:]...)
	registry.modules = append(builtinModules, applicationModules...)
	return nil
}
