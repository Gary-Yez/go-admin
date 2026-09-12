package sys_role

import (
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"
	"time"
)

type SysRole struct {
	DefaultMenu  string              `json:"default_menu" gorm:"size:100;default:sys_home;comment:默认菜单Key"`
	Id           uint                `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt    time.Time           `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt    time.Time           `json:"updated_at" gorm:"comment:更新时间"`
	Name         string              `json:"name" gorm:"unique;comment:角色名"`
	IsSuperAdmin bool                `json:"is_super_admin"`
	Menus        []*sys_menu.SysMenu `json:"menus" gorm:"many2many:sys_role_menu;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Apis         []*SysCasbinApi     `json:"apis" gorm:"-"`
}

type SysCasbinApi struct {
	Method string `json:"method" binding:"required"`
	Path   string `json:"path" binding:"required"`
}

type ApiOption struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Group       string `json:"group"`
	Description string `json:"description"`
}

type CopyBody struct {
	Id   uint   `json:"id" binding:"required"`
	Name string `json:"name" binding:"required,max=100"`
}
