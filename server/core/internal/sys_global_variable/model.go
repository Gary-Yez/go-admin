package sys_global_variable

import (
	"gitee.com/mxcker/go-admin/server/core/models/common"
)

type SysGlobalVariable struct {
	common.DbBaseModel
	Name   string `json:"name" binding:"required"`
	Type   string `json:"type" binding:"required"`
	Key    string `json:"key" gorm:"unique;" binding:"required"`
	Value  string `json:"value"`
	Remark string `json:"remark"`
}
