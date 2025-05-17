package types

import (
	"gorm.io/gorm"
	"time"
)

type DbBaseModel struct {
	Id        uint      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:更新时间"`
}

type DbBaseModelWidthDelete struct {
	Id        uint           `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}

type AuthUser struct {
	UserId uint `json:"user_id"`
	RoleId uint `json:"role_id"`
}
