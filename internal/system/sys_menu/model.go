package sys_menu

import (
	"time"
)

type SortBody struct {
	Id             uint   `json:"id" binding:"required,gt=0"`
	TargetId       uint   `json:"target_id" binding:"required,gt=0"`
	Position       string `json:"position" binding:"required,oneof=before after inside"`
	SourceParentId *uint  `json:"source_parent_id"`
	TargetParentId *uint  `json:"target_parent_id"`
	SourceIds      []uint `json:"source_ids" binding:"required,min=1,dive,gt=0"`
	TargetIds      []uint `json:"target_ids" binding:"omitempty,dive,gt=0"`
}

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
