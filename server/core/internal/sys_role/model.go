package sys_role

import (
	"gitee.com/mxcker/go-admin/server/core/common"
	"gitee.com/mxcker/go-admin/server/core/internal/sys_menu"
)

type SysRole struct {
	common.DbBaseModel
	Name    string              `json:"name" gorm:"unique;comment:角色名"`
	Default bool                `json:"default"`
	Menus   []*sys_menu.SysMenu `json:"menus" gorm:"many2many:sys_role_menu;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
