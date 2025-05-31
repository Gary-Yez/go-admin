package sys_role

import (
	"gitee.com/mxcker/go-admin/server/global"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/gin-gonic/gin"
	"time"
)

var controller = new(controllerStruct)
var Service = new(serviceStruck)
var Enforcer *casbin.SyncedCachedEnforcer

type Mounter struct{}

func (_ *Mounter) AdminRouter(adminAuthGroup *gin.RouterGroup) {
	adminAuthGroup.GET("get", controller.Get)
	adminAuthGroup.GET("list", controller.List)
	adminAuthGroup.POST("create", controller.Create)
	adminAuthGroup.POST("delete", controller.Delete)
	adminAuthGroup.POST("edit", controller.Edit)
	adminAuthGroup.POST("permission", controller.UpdatePermission)
}

func (_ *Mounter) PublicRouter(publicGroup *gin.RouterGroup) {}

func (_ *Mounter) Initialize() error {
	err := global.DB.AutoMigrate(SysRole{})
	if err != nil {
		return err
	}
	a, err := gormadapter.NewAdapterByDBUseTableName(global.DB, "sys_", "")
	if err != nil {
		return err
	}
	text := `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = r.sub == p.sub && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")
		`
	m, err := model.NewModelFromString(text)
	if err != nil {
		return err
	}
	Enforcer, _ = casbin.NewSyncedCachedEnforcer(m, a)
	if global.Config.Redis.IsNotEmpty() {
		//Enforcer.EnableCache(false)
		watcher, err := rediswatcher.NewWatcher(global.Config.Redis.Address(), rediswatcher.WatcherOptions{
			Options:                *global.Config.Redis.Option(),
			IgnoreSelf:             true,
			OptionalUpdateCallback: rediswatcher.DefaultUpdateCallback(Enforcer),
		})
		if err != nil {
			return err
		}
		err = Enforcer.SetWatcher(watcher)
		if err != nil {
			return err
		}
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
