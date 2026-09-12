package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// 保留查询链式调用方式，非法参数通过 GORM Error 传回业务层。
func queryError(db *gorm.DB, err error) *gorm.DB {
	db = db.Session(&gorm.Session{})
	_ = db.AddError(err)
	return db
}

func (req *ReqList) Validate() error {
	if req == nil {
		return errors.New("查询参数不能为空")
	}
	for i, filter := range req.Filters {
		if strings.TrimSpace(filter.Field) == "" {
			return fmt.Errorf("第%d个筛选条件缺少字段", i+1)
		}
		valid := false
		switch filter.Operator {
		case "=", "!=", ">", "<", ">=", "<=":
			valid = filterScalar(filter.Value)
		case "like":
			_, valid = filter.Value.(string)
		case "in", "between":
			values, ok := filterValues(filter.Value)
			valid = ok && len(values) > 0 && (filter.Operator != "between" || len(values) == 2)
			for _, value := range values {
				valid = valid && filterScalar(value)
			}
		default:
			return fmt.Errorf("第%d个筛选条件的操作符无效", i+1)
		}
		if !valid {
			return fmt.Errorf("第%d个筛选条件的值与操作符 %s 不匹配", i+1, filter.Operator)
		}
	}
	for i, sort := range req.Sorts {
		if strings.TrimSpace(sort.Field) == "" {
			return fmt.Errorf("第%d个排序条件缺少字段", i+1)
		}
		if order := strings.ToLower(sort.Order); order != "asc" && order != "desc" {
			return fmt.Errorf("第%d个排序方向仅支持 asc 或 desc", i+1)
		}
	}
	return nil
}

func filterScalar(value any) bool {
	if value == nil {
		return false
	}
	switch value := value.(type) {
	case time.Time:
		return true
	case json.Number:
		number, err := strconv.ParseFloat(string(value), 64)
		return err == nil && !math.IsNaN(number) && !math.IsInf(number, 0)
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return !math.IsNaN(v.Float()) && !math.IsInf(v.Float(), 0)
	}
	return false
}

func filterValues(value any) ([]any, bool) {
	v := reflect.ValueOf(value)
	if !v.IsValid() || (v.Kind() != reflect.Slice && v.Kind() != reflect.Array) {
		return nil, false
	}
	values := make([]any, v.Len())
	for i := range values {
		values[i] = v.Index(i).Interface()
	}
	return values, true
}
