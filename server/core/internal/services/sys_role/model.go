package sys_role

import (
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_menu"
	"time"
)

type SysRole struct {
	Id        uint                `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time           `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time           `json:"updated_at" gorm:"comment:更新时间"`
	Name      string              `json:"name" gorm:"unique;comment:角色名"`
	Default   bool                `json:"default"`
	Menus     []*sys_menu.SysMenu `json:"menus" gorm:"many2many:sys_role_menu;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Apis      []*SysCasbinApi     `json:"apis" gorm:"-"`
}

type SysCasbinApi struct {
	Method string `json:"method" binding:"required"`
	Path   string `json:"path" binding:"required"`
}
