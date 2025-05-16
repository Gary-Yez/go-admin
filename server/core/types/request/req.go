package request

import "gorm.io/gorm"

type Req struct {
	Id uint `json:"id" form:"id" binding:"required"`
}

func (r *Req) WithQuery(db *gorm.DB) *gorm.DB {
	return db.Where("id = ?", r.Id)
}
