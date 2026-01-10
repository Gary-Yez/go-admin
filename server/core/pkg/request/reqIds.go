package request

import (
	"gorm.io/gorm"
)

type ReqIds struct {
	Ids []uint `json:"ids" form:"ids" binding:"required"`
}

func (r *ReqIds) WithQuery(db *gorm.DB) *gorm.DB {
	return db.Where("id IN ?", r.Ids)
}
