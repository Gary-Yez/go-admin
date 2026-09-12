package sys_apis

import (
	"github.com/Gary-Yez/go-admin/request"
	"time"
)

type ApiListQuery struct {
	request.ReqList
	Keyword string `json:"keyword" form:"keyword" binding:"max=100"`
	Method  string `json:"method" form:"method" binding:"omitempty,oneof=GET POST PUT PATCH DELETE HEAD OPTIONS CONNECT TRACE"`
	Group   string `json:"group" form:"group" binding:"max=200"`
}

type SysApi struct {
	Id          uint      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt   time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"comment:更新时间"`
	Method      string    `json:"method" binding:"required" gorm:"index;uniqueIndex:idx_sys_api_method_path"`
	Path        string    `json:"path" binding:"required" gorm:"index;uniqueIndex:idx_sys_api_method_path"`
	Group       string    `json:"group" gorm:"index"`
	Description string    `json:"description"`
}

// API 方法和路径由路由维护，管理页面仅编辑展示信息。
type ApiEditBody struct {
	Id          uint   `json:"id" binding:"required"`
	Group       string `json:"group" binding:"required,max=200"`
	Description string `json:"description" binding:"required,max=255"`
}
