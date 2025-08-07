package request

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type Filter struct {
	// 操作符字段
	Operator string `json:"operator" form:"operator" binding:"required"`
	Field    string `json:"field" form:"field" binding:"required"`
	Value    any    `json:"value" form:"value"`
}

type Sort struct {
	Field string `json:"field" form:"field" binding:"required"`
	Order string `json:"order" form:"order" binding:"required"` // 枚举校验
}

type ReqList struct {
	Page    int      `json:"page" form:"page"`
	Limit   int      `json:"limit" form:"limit"`
	Filters []Filter `json:"filters" form:"filters"`
	Sorts   []Sort   `json:"sorts" form:"sorts"`
}

func (reqL *ReqList) WithFilter(db *gorm.DB, allowFields []string) *gorm.DB {
	if len(reqL.Filters) != 0 && len(allowFields) != 0 {
		// 构建字段白名单快速查找表
		allowMap := make(map[string]bool, len(allowFields))
		for _, f := range allowFields {
			allowMap[f] = true
		}
		for _, filter := range reqL.Filters {
			// 白名单校验
			if _, ok := allowMap[filter.Field]; !ok {
				continue
			}
			switch filter.Operator {
			case "=", "!=", ">", "<", ">=", "<=":
				db = db.Where(fmt.Sprintf("%s %s ?", filter.Field, filter.Operator), filter.Value)
				break
			case "like":
				if v, ok := filter.Value.(string); ok {
					//// 可选方案2：禁止用户传入通配符（转义_和%）
					//value := strings.ReplaceAll(v, "_", "\\_")
					//value = strings.ReplaceAll(value, "%", "\\%")
					esc := strings.ReplaceAll(v, "_", "\\_")
					esc = strings.ReplaceAll(esc, "%", "\\%")
					db = db.Where(fmt.Sprintf("%s LIKE ?", filter.Field), "%"+esc+"%")
					//db = db.Where(fmt.Sprintf("%s LIKE ?", filter.Field), v)
				}
				break
			case "in":
				// expect filter.Value to be a slice, e.g. []interface{}{"a","b"}
				db = db.Where(fmt.Sprintf("%s IN (?)", filter.Field), filter.Value)
			case "between":
				if vals, ok := filter.Value.([]interface{}); ok && len(vals) == 2 {
					db = db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", filter.Field), vals[0], vals[1])
				}
			default:
				continue
			}
		}
	}
	return db
}

func (reqL *ReqList) WithSort(db *gorm.DB, allowFields []string) *gorm.DB {
	if len(allowFields) != 0 {
		allowMap := make(map[string]bool, len(allowFields))
		for _, f := range allowFields {
			allowMap[f] = true
		}
		// 排序逻辑
		for _, sort := range reqL.Sorts {
			// 白名单校验
			if _, ok := allowMap[sort.Field]; !ok {
				continue
			}
			if strings.ToLower(sort.Order) == "desc" {
				db = db.Order(fmt.Sprintf("%s DESC", sort.Field))
			} else {
				db = db.Order(fmt.Sprintf("%s ASC", sort.Field))
			}
		}
	}
	return db
}

func (reqL *ReqList) WithPagination(db *gorm.DB) *gorm.DB {
	if reqL.Page > 0 {
		if reqL.Limit == 0 {
			reqL.Limit = 10
		}
		if reqL.Limit > 100 {
			reqL.Limit = 100
		}
		limit := reqL.Limit
		offset := reqL.Limit * (reqL.Page - 1)
		db = db.Limit(limit).Offset(offset)
	}
	return db
}
