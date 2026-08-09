package sys_apis

import (
	"errors"
	"gitee.com/mxcker/go-admin/internal/state"

	"gitee.com/mxcker/go-admin/internal/system/sys_role"
	request2 "gitee.com/mxcker/go-admin/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type serviceStruct struct {
}

func (s *serviceStruct) SyncApi() (newApis, deleteApis, ignoreApis []*SysApi, err error) {
	newApis = make([]*SysApi, 0)
	deleteApis = make([]*SysApi, 0)
	ignoreApis = make([]*SysApi, 0)
	var apis []*SysApi
	if err = state.DB().Find(&apis).Error; err != nil {
		return
	}
	var ignores []*SysIgnoreApi
	if err = state.DB().Find(&ignores).Error; err != nil {
		return
	}
	// 数据库中所有的API Map
	var apisMap = make(map[string]bool)
	for _, api := range apis {
		key := api.Method + "_" + api.Path
		apisMap[key] = true
	}
	// 数据库中所有忽略的API Map
	var ignoresMap = make(map[string]bool)
	for _, ignore := range ignores {
		key := ignore.Method + "_" + ignore.Path
		ignoresMap[key] = true
	}
	// 实际注册的路由API map
	var routeMap = make(map[string]bool)
	for _, route := range state.Routes() {
		path := strings.TrimPrefix(route.Path, state.Config().Server.ApiPrefix)
		key := route.Method + "_" + path
		routeMap[key] = true
		// api被忽略则跳出循环
		if ignoresMap[key] {
			continue
		}
		// 判断需要添加的新API
		if !apisMap[key] {
			newApis = append(newApis, &SysApi{
				Group:       "",
				Description: "",
				Method:      route.Method,
				Path:        path,
			})
		}
	}
	//需要删除的API
	for _, api := range apis {
		if !routeMap[api.Method+"_"+api.Path] {
			deleteApis = append(deleteApis, api)
		}
	}
	// 需要删除的忽略API
	needDeleteIgnoreIds := make([]uint, 0)
	for _, ignore := range ignores {
		if !routeMap[ignore.Method+"_"+ignore.Path] {
			//需要删除的忽略API
			needDeleteIgnoreIds = append(needDeleteIgnoreIds, ignore.Id)
		} else {
			ignoreApis = append(ignoreApis, &SysApi{
				Method:      ignore.Method,
				Path:        ignore.Path,
				Group:       "",
				Description: "",
			})
		}
	}
	if len(needDeleteIgnoreIds) > 0 {
		err = state.DB().Delete(&SysIgnoreApi{}, needDeleteIgnoreIds).Error
	}
	return
}

func (s *serviceStruct) GetGroups() (groups []string, err error) {
	groups = make([]string, 0)
	err = state.DB().Model(&SysApi{}).Where("`group` != ''").Distinct("`group`").Pluck("`group`", &groups).Error
	return
}

func (s *serviceStruct) Get(req *request2.Req) (data *SysApi, err error) {
	data = &SysApi{}
	err = req.WithQuery(state.DB().Model(SysApi{})).First(data).Error
	return
}

func (s *serviceStruct) List(req *request2.ReqList) (list []*SysApi, total int64, err error) {
	db := req.WithFilter(state.DB().Model(SysApi{}), nil)
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = req.WithPagination(req.WithSort(db, nil)).Find(&list).Error
	return
}

func (s *serviceStruct) CreateOrUpdate(data *SysApi) (err error) {
	return state.DB().Create(data).Error
}

func (s *serviceStruct) Update(data *SysApi) (err error) {
	if data.Id == 0 {
		return errors.New("id不能为空")
	}
	return state.DB().Select("*").
		Omit(clause.Associations).
		Omit("Id", "CreatedAt", "UpdatedAt").
		Where("id = ?", data.Id).Updates(data).Error
}

func (s *serviceStruct) UpdateIgnore(api *SysIgnoreApi) (err error) {
	if api.Ignore {
		return state.DB().Create(api).Error
	} else {
		return state.DB().Where("`path` = ? AND method = ?", api.Path, api.Method).Delete(api).Error
	}
}

func (s *serviceStruct) DeleteByIds(req *request2.ReqIds) error {
	return state.DB().Transaction(func(tx *gorm.DB) error {
		var list []*SysApi
		if err := req.WithQuery(state.DB()).Find(&list).Error; err != nil {
			return err
		}
		if err := req.WithQuery(state.DB()).Delete(&SysApi{}).Error; err != nil {
			return err
		}
		for _, sysApi := range list {
			err := sys_role.Service.DeletePolicy(1, sysApi.Path, sysApi.Method)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
