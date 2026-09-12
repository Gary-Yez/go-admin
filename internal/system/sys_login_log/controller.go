package sys_login_log

import (
	"github.com/Gary-Yez/go-admin/request"
	"github.com/Gary-Yez/go-admin/response"
	"github.com/gin-gonic/gin"
)

type controllerStruct struct{}

func (*controllerStruct) List(ctx *gin.Context) {
	req, err := request.GetReqList(ctx)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	list, total, err := Service.List(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.List(ctx, list, total)
}
func (*controllerStruct) Delete(ctx *gin.Context) {
	req, err := request.GetReqIds(ctx)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	if err := Service.Delete(ctx.Request.Context(), req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "登录日志已删除")
}
func (*controllerStruct) Cleanup(ctx *gin.Context) {
	body := new(CleanupBody)
	if err := ctx.ShouldBindJSON(body); err != nil {
		response.Error(ctx, err)
		return
	}
	count, err := Service.Cleanup(ctx.Request.Context(), body.Days)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"count": count})
}
