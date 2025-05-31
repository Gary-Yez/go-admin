package sys_apis

import "gitee.com/mxcker/go-admin/server/types"

type SysApi struct {
	types.DbBaseModel
	Method      string `json:"method" binding:"required" gorm:"index"`
	Path        string `json:"path" binding:"required" gorm:"index"`
	Group       string `json:"group" gorm:"index"`
	Description string `json:"description"`
}

type SysIgnoreApi struct {
	types.DbBaseModel
	Method string `json:"method" binding:"required"`
	Path   string `json:"path" binding:"required"`
	Ignore bool   `json:"ignore" gorm:"-"`
}
