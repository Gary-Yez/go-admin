package sys_role

import (
	"errors"
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/pkg/request"
	"gitee.com/mxcker/go-admin/server/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

type serviceStruck struct {
}

func (s *serviceStruck) Get(req *request.Req) (data *SysRole, err error) {
	data = &SysRole{}
	err = req.WithQuery(global.DB.Model(SysRole{}).Preload("Menus")).First(data).Error
	policy, err := Enforcer.GetFilteredPolicy(0, strconv.Itoa(int(data.Id)))
	if err != nil {
		return nil, err
	}
	pathMaps := make([]*SysCasbinApi, 0)
	for _, v := range policy {
		pathMaps = append(pathMaps, &SysCasbinApi{
			Path:   v[1],
			Method: v[2],
		})
	}
	data.Apis = pathMaps
	return
}

func (s *serviceStruck) List() (list []*SysRole, total int64, err error) {
	err = global.DB.Model(SysRole{}).Find(&list).Error
	return
}

func (s *serviceStruck) Create(data *SysRole) (err error) {
	err = global.DB.Create(data).Error
	return
}

func (s *serviceStruck) Update(data *SysRole) (err error) {
	if data.Id == 0 {
		return errors.New("id不能为空")
	}
	return global.DB.Select("*").
		Omit(clause.Associations).
		Omit("Id", "CreatedAt", "UpdatedAt").
		Where("id = ?", data.Id).Updates(data).Error
}

func (s *serviceStruck) DeleteByIds(req *request.ReqIds) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		err := req.WithQuery(global.DB).Where("`default` = ?", false).Delete(&SysRole{}).Error
		if err != nil {
			if utils.IsForeignKeyConstraintError(err) {
				return errors.New("该角色正在使用中")
			}
			return err
		}
		for _, v := range req.Ids {
			err = s.DeletePolicy(0, strconv.Itoa(int(v)))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *serviceStruck) UpdatePermission(data *SysRole) (err error) {
	roleId := strconv.Itoa(int(data.Id))
	err = s.DeletePolicy(0, roleId)
	if err != nil {
		return err
	}
	var rules [][]string
	//做权限去重处理
	deduplicateMap := make(map[string]bool)
	for _, v := range data.Apis {
		key := roleId + v.Path + v.Method
		if _, ok := deduplicateMap[key]; !ok {
			deduplicateMap[key] = true
			rules = append(rules, []string{roleId, v.Path, v.Method})
		}
	}
	if len(rules) != 0 {
		success, err := Enforcer.AddPolicies(rules)
		if err != nil {
			return err
		}
		if !success {
			return errors.New("API权限修改失败")
		}
	} // 设置空权限无需调用 AddPolicies 方法
	return global.DB.Model(data).Association("Menus").Replace(data.Menus)
}

func (s *serviceStruck) DeletePolicy(fieldIndex int, fieldValues ...string) error {
	_, err := Enforcer.RemoveFilteredPolicy(fieldIndex, fieldValues...)
	return err
}
