package sys_task

import "gitee.com/mxcker/go-admin/server/core/types"

type SysTask struct {
	types.DbBaseModel
	Name             string `json:"name" binding:"required"`
	HandlerKey       string `json:"handler_key" binding:"required"`
	Params           string `json:"params"`
	Cron             string `json:"cron" binding:"required"`
	AllowConcurrency bool   `json:"allow_concurrency"`
	Enable           bool   `json:"enable"`
}
