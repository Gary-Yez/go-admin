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
	"github.com/Gary-Yez/go-admin/internal/system/sys_api_token"
	"github.com/Gary-Yez/go-admin/internal/system/sys_apis"
	"github.com/Gary-Yez/go-admin/internal/system/sys_auth"
	"github.com/Gary-Yez/go-admin/internal/system/sys_config"
	"github.com/Gary-Yez/go-admin/internal/system/sys_cron_job"
	"github.com/Gary-Yez/go-admin/internal/system/sys_devtools"
	"github.com/Gary-Yez/go-admin/internal/system/sys_file"
	"github.com/Gary-Yez/go-admin/internal/system/sys_login_log"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"
	"github.com/Gary-Yez/go-admin/internal/system/sys_monitor"
	"github.com/Gary-Yez/go-admin/internal/system/sys_role"
	"github.com/Gary-Yez/go-admin/internal/system/sys_storage"
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

var engineCallbacks []func(*gin.Engine) error
var engineConfigurationClosed bool

// ConfigureEngine 在 Run 前注册配置，按顺序在框架路由注册前执行。
// 回调直接注册的路由不参与框架认证和 API 同步；业务路由应通过 Module 注册。
func ConfigureEngine(configure func(*gin.Engine) error) error {
	if engineConfigurationClosed {
		return errors.New("请在 admin.Run 之前注册 Gin 配置")
	}
	if configure == nil {
		return errors.New("Gin 配置回调不能为空")
	}
	engineCallbacks = append(engineCallbacks, configure)
	return nil
}

func createEngine() (*gin.Engine, error) {
	engine := gin.Default()
	for i, configure := range engineCallbacks {
		if err := configure(engine); err != nil {
			return nil, fmt.Errorf("执行第 %d 个 Gin 配置回调失败：%w", i+1, err)
		}
	}
	return engine, nil
}

// Run initializes the singleton admin runtime and starts its HTTP server.
// 默认读取运行目录下的 config.yaml，命令行 --config/-c 优先于传入路径。
func Run(configPaths ...string) error {
	if engineConfigurationClosed {
		return errors.New("admin.Run 不能重复调用")
	}
	engineConfigurationClosed = true
	if len(configPaths) > 1 {
		return errors.New("admin.Run只能指定一个配置文件路径")
	}
	configFile := "config.yaml"
	if len(configPaths) == 1 && strings.TrimSpace(configPaths[0]) != "" {
		configFile = configPaths[0]
	}
	flags := pflag.NewFlagSet("go-admin", pflag.ContinueOnError)
	flags.StringVarP(&configFile, "config", "c", configFile, "配置文件路径（YAML）")
	flags.String("server.host", "0.0.0.0", "Web服务运行的IP")
	flags.String("server.port", "8080", "Web服务运行的端口")
	if err := flags.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			return nil
		}
		return fmt.Errorf("parse arguments: %w", err)
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("不支持位置参数，请使用 --config 或 -c 指定配置文件")
	}
	cfg, err := initialization.InitConfig(flags, configFile)
	if err != nil {
		return err
	}

	deps, err := initialization.InitDependencies(cfg)
	if err != nil {
		return err
	}
	state.Configure(cfg, deps.DB, deps.Cache, deps.Scheduler)

	if err := registerBuiltins(); err != nil {
		return err
	}
	configSchema, err = initialization.InitBusinessConfig[BaseConfig](configSchema)
	if err != nil {
		return fmt.Errorf("初始化业务配置：%w", err)
	}
	configReady = true
	for _, module := range registry.modules {
		log.Println("初始化模块：", module.Name())
		if err := module.Initialize(); err != nil {
			return fmt.Errorf("initialize module %q: %w", module.key, err)
		}
	}
	var menus []MenuDefinition
	for _, module := range registry.modules {
		if provider, ok := module.Module.(MenuProvider); ok {
			menus = append(menus, provider.Menus()...)
		}
	}
	if len(menus) > 0 {
		if err := sys_menu.Sync(menus); err != nil {
			return fmt.Errorf("初始化业务菜单：%w", err)
		}
	}
	registry.initialized = true

	if cfg.IsDev() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine, err := createEngine()
	if err != nil {
		return err
	}
	engine.Use(static.Serve(cfg.Server.AdminPrefix, static.LocalFile("./dist", false)))

	apiPrefix := strings.TrimSpace(cfg.Server.ApiPrefix)
	adminGroup := engine.Group(apiPrefix,
		middlewares.JWTMiddleware(sys_api_token.Service.VerifyApiToken, sys_auth.Service.VerifyAuthUser),
		middlewares.CasbinMiddleware(apiPrefix, sys_role.Enforcer),
	)
	publicGroup := engine.Group(apiPrefix)
	// 回调可能已注册公共接口，只收集随后由 AdminRouter 新增的路由。
	existingRoutes := make(map[string]struct{})
	for _, route := range engine.Routes() {
		existingRoutes[route.Method+" "+route.Path] = struct{}{}
	}
	moduleNames := make(map[string]string, len(registry.modules))
	for _, module := range registry.modules {
		moduleNames[module.key] = module.Name()
		module.AdminRouter(adminGroup.Group(module.key))
	}
	// 在注册公共路由前记录受保护路由，API 同步只使用 AdminRouter 的接口。
	routes := gin.RoutesInfo{}
	for _, route := range engine.Routes() {
		if _, exists := existingRoutes[route.Method+" "+route.Path]; !exists {
			routes = append(routes, route)
		}
	}
	slices.SortFunc(routes, func(left, right gin.RouteInfo) int {
		return cmp.Compare(left.Path, right.Path)
	})
	state.SetAdminRoutes(routes)
	for _, module := range registry.modules {
		module.PublicRouter(publicGroup.Group(module.key))
	}
	if err := sys_apis.Service.SyncNewAPIs(moduleNames); err != nil {
		return fmt.Errorf("自动同步新增API：%w", err)
	}
	stopMonitor := initialization.StartMonitor()
	defer stopMonitor()
	state.Scheduler().StartScheduler(time.Second)
	defer func() {
		if err := state.Scheduler().StopScheduler(); err != nil {
			log.Printf("停止定时任务失败：%v", err)
		}
	}()
	stopPolicyRefresh := sys_role.StartPolicyRefresh()
	defer stopPolicyRefresh()

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
		{key: "sys_login_log", Module: new(sys_login_log.Mounter)},
		{key: "sys_monitor", Module: new(sys_monitor.Mounter)},
		{key: "sys_storage", Module: new(sys_storage.Mounter)},
		{key: "sys_file", Module: new(sys_file.Mounter)},
		{key: "sys_auth", Module: new(sys_auth.Mounter)},
		{key: "sys_api_token", Module: new(sys_api_token.Mounter)},
		{key: "sys_cron_job", Module: new(sys_cron_job.Mounter)},
		{key: "sys_config", Module: new(sys_config.Mounter)},
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
