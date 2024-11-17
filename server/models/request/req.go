package request

import "gorm.io/gorm"

type Req struct {
	Id uint `json:"id" form:"id" binding:"required"`
}

type ReqIds struct {
	Ids []uint `json:"ids" form:"ids" binding:"required"`
}

type ReqList struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

func (reqL *ReqList) SetDB(db *gorm.DB) {
	if reqL.Page > 0 {
		if reqL.Limit == 0 {
			reqL.Limit = 10
		}
		limit := reqL.Limit
		offset := reqL.Limit * (reqL.Page - 1)
		db = db.Limit(limit).Offset(offset)
	}
}
