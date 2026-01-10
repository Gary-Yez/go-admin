package test

import (
	"gorm.io/gorm"
	"time"
)

type Test struct {
    Id        uint      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
    CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
    UpdatedAt time.Time `json:"updated_at" gorm:"comment:更新时间"`
    DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;comment:删除时间"`
	A string `json:"a" binding:"required"`
}

