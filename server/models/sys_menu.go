package models

type SysMenu struct {
	DbBaseModel
	Name      string     `json:"name"`
	Icon      string     `json:"icon"`
	Path      string     `json:"path" gorm:"unique"`
	Component string     `json:"component"`
	Sort      int        `json:"sort"`
	Hidden    bool       `json:"hidden"`
	Children  []*SysMenu `json:"children" gorm:"foreignKey:ParentId"`
	ParentId  *uint      `json:"parent_id" gorm:"default:null"`
}
