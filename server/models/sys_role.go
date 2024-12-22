package models

type SysRole struct {
	DbBaseModel
	Name    string     `json:"name" gorm:"unique;comment:角色名"`
	Default bool       `json:"default"`
	Menus   []*SysMenu `json:"menus" gorm:"many2many:sys_role_menu;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
