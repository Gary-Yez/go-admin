package sys_apis

import (
	"time"
)

type SysApi struct {
	Id          uint      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt   time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"comment:更新时间"`
	Method      string    `json:"method" binding:"required" gorm:"index"`
	Path        string    `json:"path" binding:"required" gorm:"index"`
	Group       string    `json:"group" gorm:"index"`
	Description string    `json:"description"`
}

type SysIgnoreApi struct {
	Id        uint      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:更新时间"`
	Method    string    `json:"method" binding:"required"`
	Path      string    `json:"path" binding:"required"`
	Ignore    bool      `json:"ignore" gorm:"-"`
}
