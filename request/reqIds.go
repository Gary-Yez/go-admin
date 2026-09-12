package request

import (
	"errors"
	"gorm.io/gorm"
)

type ReqIds struct {
	Ids []uint `json:"ids" form:"ids" binding:"required,min=1,dive,gt=0"`
}

func (r *ReqIds) Validate() error {
	if r == nil || len(r.Ids) == 0 {
		return errors.New("批量操作需要提供ID")
	}
	seen := make(map[uint]bool, len(r.Ids))
	ids := make([]uint, 0, len(r.Ids))
	for _, id := range r.Ids {
		if id == 0 {
			return errors.New("ID必须大于0")
		}
		if !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	r.Ids = ids
	return nil
}

func (r *ReqIds) WithQuery(db *gorm.DB) *gorm.DB {
	if err := r.Validate(); err != nil {
		return queryError(db, err)
	}
	return db.Where("id IN ?", r.Ids)
}
