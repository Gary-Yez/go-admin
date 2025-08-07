package sys_role

import (
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_menu"
	"gitee.com/mxcker/go-admin/server/utils"
)

type SysRole struct {
	utils.DbBaseModel
	Name    string              `json:"name" gorm:"unique;comment:角色名"`
	Default bool                `json:"default"`
	Menus   []*sys_menu.SysMenu `json:"menus" gorm:"many2many:sys_role_menu;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Apis    []*SysCasbinApi     `json:"apis" gorm:"-"`
}

type SysCasbinApi struct {
	Method string `json:"method" binding:"required"`
	Path   string `json:"path" binding:"required"`
}
