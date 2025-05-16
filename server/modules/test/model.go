package test

import "gitee.com/mxcker/go-admin/server/core/types"

type Test struct {
	types.DbBaseModel
	Test string `json:"test" binding:"required"`
	Age int `json:"age" binding:"required"`
}

