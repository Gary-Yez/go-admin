package sys_role

import (
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"
	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	"github.com/casbin/casbin/v3/persist"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"time"
)

var controller = new(controllerStruct)
var Service = new(serviceStruck)
var Enforcer *casbin.SyncedCachedEnforcer
var policyWatcher persist.Watcher

type Mounter struct{}

func (_ *Mounter) Name() string {
	return "核心服务-角色管理"
}

func (_ *Mounter) Menus() []sys_menu.Definition {
	return []sys_menu.Definition{
		{
			Name: "角色管理", Key: "sys_role", ParentKey: "sys_permission", Icon: "iconoir:shield-check",
			Path: "sys_role", Component: "../core/views/sys_role/index.vue", Sort: 2,
		},
	}
}

func (_ *Mounter) AdminRouter(adminAuthGroup *gin.RouterGroup) {
	adminAuthGroup.GET("get", controller.Get)
	adminAuthGroup.GET("list", controller.List)
	adminAuthGroup.POST("create", controller.Create)
	adminAuthGroup.POST("copy", controller.Copy)
	adminAuthGroup.POST("delete", controller.Delete)
	adminAuthGroup.POST("edit", controller.Edit)
	adminAuthGroup.POST("permission", controller.UpdatePermission)
}

func (_ *Mounter) PublicRouter(publicGroup *gin.RouterGroup) {}

func (_ *Mounter) Initialize() error {
	// 重命名数据库列，保留现有超级管理员标识。
	migrator := state.DB().Migrator()
	if migrator.HasColumn(&SysRole{}, "default") && !migrator.HasColumn(&SysRole{}, "is_super_admin") {
		if err := migrator.RenameColumn(&SysRole{}, "default", "is_super_admin"); err != nil {
			return err
		}
	}
	err := state.DB().AutoMigrate(SysRole{})
	if err != nil {
		return err
	}
	a, err := gormadapter.NewAdapterByDBUseTableName(state.DB(), "sys_", "")
	if err != nil {
		return err
	}
	text := `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = r.sub == p.sub && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")
		`
	m, err := model.NewModelFromString(text)
	if err != nil {
		return err
	}
	Enforcer, err = casbin.NewSyncedCachedEnforcer(m, a)
	if err != nil {
		return err
	}
	policyWatcher = nil
	if state.Config().Redis.IsNotEmpty() {
		watcher, err := newPolicyWatcher(Enforcer, state.Config())
		if err != nil {
			return err
		}
		policyWatcher = watcher
	}
	Enforcer.SetExpireTime(time.Hour)
	err = Enforcer.LoadPolicy()
	if err != nil {
		return err
	}
	err = InitData()
	if err != nil {
		return err
	}
	return nil
}
