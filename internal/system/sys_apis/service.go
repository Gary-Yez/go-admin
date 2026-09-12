package sys_apis

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gary-Yez/go-admin/internal/permissions"
	"github.com/Gary-Yez/go-admin/internal/state"

	"github.com/Gary-Yez/go-admin/internal/system/sys_role"
	request2 "github.com/Gary-Yez/go-admin/request"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type serviceStruct struct {
}

func (s *serviceStruct) GetInvalidAPIs() ([]*SysApi, error) {
	_, apis, err := s.syncApi(state.DB())
	return apis, err
}

// 启动时只补录新增接口，不覆盖已有信息、不删除接口或修改角色授权。
func (s *serviceStruct) SyncNewAPIs(moduleNames map[string]string) error {
	return state.DB().Transaction(func(tx *gorm.DB) error {
		newApis, _, err := s.syncApi(tx)
		if err != nil {
			return err
		}
		if len(newApis) == 0 {
			return nil
		}
		defaults := make(map[string]SysApi)
		for _, api := range defaultAPIs() {
			defaults[api.Method+" "+api.Path] = api
		}
		for _, api := range newApis {
			if definition, exists := defaults[api.Method+" "+api.Path]; exists {
				api.Group, api.Description = definition.Group, definition.Description
				continue
			}
			module, _, _ := strings.Cut(strings.TrimPrefix(api.Path, "/"), "/")
			name := strings.TrimSpace(moduleNames[module])
			if name == "" {
				name = module
			}
			api.Group = "业务模块-" + name
			if strings.HasPrefix(module, "sys_") {
				api.Group = name
			}
			api.Description = api.Method + " " + api.Path
		}
		// 多实例同时启动时，由方法和路径的唯一索引避免重复录入。
		return tx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(newApis, 100).Error
	})
}

func (s *serviceStruct) syncApi(db *gorm.DB) (newApis, deleteApis []*SysApi, err error) {
	newApis = make([]*SysApi, 0)
	deleteApis = make([]*SysApi, 0)
	var apis []*SysApi
	if err = db.Find(&apis).Error; err != nil {
		return
	}
	// 数据库中所有的API Map
	var apisMap = make(map[string]bool)
	for _, api := range apis {
		key := api.Method + "_" + api.Path
		apisMap[key] = true
	}
	// 只同步 AdminRouter 注册的路由，公共接口不进入同步清单。
	var routeMap = make(map[string]bool)
	for _, route := range state.AdminRoutes() {
		path := strings.TrimPrefix(route.Path, state.Config().Server.ApiPrefix)
		if permissions.IsAuthenticatedAPI(route.Method, path) {
			continue
		}
		key := route.Method + "_" + path
		routeMap[key] = true
		// 判断需要添加的新API
		if !apisMap[key] {
			newApis = append(newApis, &SysApi{
				Method: route.Method,
				Path:   path,
			})
		}
	}
	//需要删除的API
	for _, api := range apis {
		if !routeMap[api.Method+"_"+api.Path] {
			deleteApis = append(deleteApis, api)
		}
	}

	return
}

func (s *serviceStruct) GetGroups() (groups []string, err error) {
	groups = make([]string, 0)
	err = state.DB().Model(&SysApi{}).Where("`group` != ''").Distinct("`group`").Pluck("`group`", &groups).Error
	return
}

func (s *serviceStruct) List(req *ApiListQuery) (list []*SysApi, total int64, err error) {
	db := req.WithFilter(state.DB().Model(SysApi{}), nil)
	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		keyword = "%" + strings.NewReplacer("!", "!!", "%", "!%", "_", "!_").Replace(keyword) + "%"
		db = db.Where("(path LIKE ? ESCAPE '!' OR description LIKE ? ESCAPE '!')", keyword, keyword)
	}
	if req.Method != "" {
		db = db.Where(clause.Eq{Column: "method", Value: req.Method})
	}
	if req.Group != "" {
		db = db.Where(clause.Eq{Column: "group", Value: req.Group})
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = req.WithPagination(req.WithSort(db, []string{"id"})).Order("id DESC").Find(&list).Error
	return
}

func (s *serviceStruct) Update(data *ApiEditBody) error {
	if data.Id == 0 {
		return errors.New("id不能为空")
	}
	return state.DB().Model(&SysApi{}).Where("id = ?", data.Id).
		Updates(map[string]interface{}{"group": data.Group, "description": data.Description}).Error
}

func (s *serviceStruct) DeleteByIds(req *request2.ReqIds) error {
	adapter, ok := sys_role.Enforcer.GetAdapter().(*gormadapter.Adapter)
	if !ok {
		return errors.New("角色权限存储未就绪")
	}
	transaction, err := adapter.BeginTransaction(context.Background())
	if err != nil {
		return err
	}
	defer transaction.Rollback()
	txAdapter := transaction.GetAdapter().(*gormadapter.Adapter)
	tx := txAdapter.GetDb().Session(&gorm.Session{NewDB: true})
	var list []*SysApi
	if err := req.WithQuery(tx).Clauses(clause.Locking{Strength: "UPDATE"}).Find(&list).Error; err != nil {
		return err
	}
	if len(list) == 0 {
		return nil
	}
	if err := req.WithQuery(tx).Delete(&SysApi{}).Error; err != nil {
		return err
	}
	for _, api := range list {
		if err := txAdapter.RemoveFilteredPolicy("p", "p", 1, api.Path, api.Method); err != nil {
			return err
		}
	}
	if err := transaction.Commit(); err != nil {
		return err
	}
	if err := sys_role.Service.ReloadPolicy(); err != nil {
		return fmt.Errorf("API及关联权限已删除，但权限缓存更新失败：%w", err)
	}
	return nil
}
