package utils

import "strings"

func IsForeignKeyConstraintError(err error) bool {
	if err == nil {
		return false
	}
	// 判断是否包含错误码 1451（父行被引用）
	return strings.Contains(err.Error(), "Error 1451")
}
