package request

import (
	"fmt"
	"gorm.io/gorm"
)

type Where struct {
	Field string `json:"field" binding:"required"`
	// 操作符字段
	Operator string `json:"operator" binding:"required,oneof=eq ne gt lt like"`
	Value    string `json:"value"`
}

type ReqList struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Where *Where `json:"where"`
}

func (reqL *ReqList) BuildQuery(db *gorm.DB) *gorm.DB {
	if reqL.Page > 0 {
		if reqL.Limit == 0 {
			reqL.Limit = 10
		}
		limit := reqL.Limit
		offset := reqL.Limit * (reqL.Page - 1)
		db = db.Limit(limit).Offset(offset)
	}
	return reqL.BuildWhere(db)
}

func (reqL *ReqList) BuildWhere(db *gorm.DB) *gorm.DB {
	if reqL.Where != nil {
		switch reqL.Where.Operator {
		case "eq":
			return db.Where(fmt.Sprintf("%s = ?", reqL.Where.Field), reqL.Where.Value)
		case "ne":
			return db.Where(fmt.Sprintf("%s != ?", reqL.Where.Field), reqL.Where.Value)
		case "gt":
			return db.Where(fmt.Sprintf("%s > ?", reqL.Where.Field), reqL.Where.Value)
		case "lt":
			return db.Where(fmt.Sprintf("%s < ?", reqL.Where.Field), reqL.Where.Value)
		case "like":
			return db.Where(fmt.Sprintf("%s LIKE ?", reqL.Where.Field), reqL.Where.Value)
		// 添加其他操作符处理
		default:
			return db
		}
	} else {
		return db
	}
}
