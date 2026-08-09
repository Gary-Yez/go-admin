package sys_menu

import "time"

type SysMenu struct {
	Id        uint       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time  `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"comment:更新时间"`
	Name      string     `json:"name"`
	Key       string     `json:"key" gorm:"unique"`
	Icon      string     `json:"icon"`
	Path      string     `json:"path"`
	Component string     `json:"component"`
	Sort      int        `json:"sort"`
	Hidden    bool       `json:"hidden"`
	Children  []*SysMenu `json:"children" gorm:"foreignKey:ParentId"`
	ParentId  *uint      `json:"parent_id" gorm:"default:null"`
}
