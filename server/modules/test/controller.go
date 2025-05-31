package test

import (
	"gitee.com/mxcker/go-admin/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type controllerStruct struct{}

func (c *controllerStruct) Test(ctx *gin.Context) {
	response.Success(ctx, "test")
}
