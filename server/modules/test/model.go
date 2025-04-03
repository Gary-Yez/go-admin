package test

import "gitee.com/mxcker/go-admin/server/core/common"

type Test struct {
	common.DbBaseModel
	A string `json:"a" binding:"required"`
}

