package sys_task

import "gitee.com/mxcker/go-admin/server/core/types"

type SysTask struct {
	types.DbBaseModel
	Name string `json:"name" binding:"required"`
}

