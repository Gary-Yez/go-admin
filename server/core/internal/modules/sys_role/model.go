package sys_role

import (
	"gitee.com/mxcker/go-admin/server/core/internal/modules/sys_menu"
	"gitee.com/mxcker/go-admin/server/core/types"
)

type SysRole struct {
	types.DbBaseModel
	Name    string              `json:"name" gorm:"unique;comment:角色名"`
	Default bool                `json:"default"`
	Menus   []*sys_menu.SysMenu `json:"menus" gorm:"many2many:sys_role_menu;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
