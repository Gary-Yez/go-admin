package request

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	Order string `json:"order" form:"order" binding:"required"`
}

// DefaultSortFields 返回默认排序字段，并追加业务允许排序的字段。
// 每次返回独立切片，修改结果不会影响其他调用。
func DefaultSortFields(extra ...string) []string {
	return append([]string{"id", "created_at", "updated_at"}, extra...)
}

type ReqList struct {
	Page    int      `json:"page" form:"page"`
	Limit   int      `json:"limit" form:"limit"`
	Filters []Filter `json:"filters" form:"filters" binding:"omitempty,dive"`
	Sorts   []Sort   `json:"sorts" form:"sorts" binding:"omitempty,dive"`
}

func (reqL *ReqList) WithFilter(db *gorm.DB, allowFields []string) *gorm.DB {
	if err := reqL.Validate(); err != nil {
		return queryError(db, err)
	}
	if len(reqL.Filters) != 0 {
		// 构建字段白名单快速查找表
		allowMap := make(map[string]bool, len(allowFields))
		for _, f := range allowFields {
			allowMap[f] = true
		}
		for _, filter := range reqL.Filters {
			// 白名单校验
			if _, ok := allowMap[filter.Field]; !ok {
				return queryError(db, fmt.Errorf("不允许筛选字段 %s", filter.Field))
			}
			switch filter.Operator {
			case "=", "!=", ">", "<", ">=", "<=":
				db = db.Where(clause.Expr{SQL: "? " + filter.Operator + " ?", Vars: []any{clause.Column{Name: filter.Field}, filter.Value}})
				break
			case "like":
				if v, ok := filter.Value.(string); ok {
					// 显式指定转义字符，让下划线、百分号按字面值匹配。
					esc := strings.NewReplacer("!", "!!", "_", "!_", "%", "!%").Replace(v)
					db = db.Where(clause.Expr{SQL: "? LIKE ? ESCAPE '!'", Vars: []any{clause.Column{Name: filter.Field}, "%" + esc + "%"}})
				}
				break
			case "in":
				// expect filter.Value to be a slice, e.g. []interface{}{"a","b"}
				db = db.Where(clause.Expr{SQL: "? IN (?)", Vars: []any{clause.Column{Name: filter.Field}, filter.Value}})
			case "between":
				if vals, ok := filterValues(filter.Value); ok && len(vals) == 2 {
					db = db.Where(clause.Expr{SQL: "? BETWEEN ? AND ?", Vars: []any{clause.Column{Name: filter.Field}, vals[0], vals[1]}})
				}
			default:
				return queryError(db, fmt.Errorf("不支持的筛选操作符 %s", filter.Operator))
			}
		}
	}
	return db
}

func (reqL *ReqList) WithSort(db *gorm.DB, allowFields []string) *gorm.DB {
	if err := reqL.Validate(); err != nil {
		return queryError(db, err)
	}
	if len(reqL.Sorts) != 0 {
		allowMap := make(map[string]bool, len(allowFields))
		for _, f := range allowFields {
			allowMap[f] = true
		}
		// 排序逻辑
		for _, sort := range reqL.Sorts {
			// 白名单校验
			if _, ok := allowMap[sort.Field]; !ok {
				return queryError(db, fmt.Errorf("不允许排序字段 %s", sort.Field))
			}
			if strings.ToLower(sort.Order) == "desc" {
				db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: sort.Field}, Desc: true})
			} else {
				db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: sort.Field}})
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
