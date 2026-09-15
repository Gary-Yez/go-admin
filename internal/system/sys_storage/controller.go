package sys_storage

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
func (*controllerStruct) Get(ctx *gin.Context) {
	var body struct {
		Id uint `form:"id" binding:"required,gt=0"`
	}
	if err := ctx.ShouldBindQuery(&body); err != nil {
		response.Error(ctx, err)
		return
	}
	result, err := Service.Get(ctx.Request.Context(), body.Id)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, result)
}
func (*controllerStruct) Save(ctx *gin.Context) {
	var body SaveBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := Service.Save(ctx.Request.Context(), &body); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx)
}
func (*controllerStruct) Delete(ctx *gin.Context) {
	req, err := request.GetReqIds(ctx)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	for _, id := range req.Ids {
		if err := Service.Delete(ctx.Request.Context(), id); err != nil {
			response.Error(ctx, err)
			return
		}
	}
	response.Success(ctx)
}

func (*controllerStruct) SetDefault(ctx *gin.Context) {
	var body struct {
		Id uint `json:"id" binding:"required,gt=0"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := Service.SetState(ctx.Request.Context(), body.Id, true, true); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx)
}
func (*controllerStruct) SetEnabled(ctx *gin.Context) {
	var body struct {
		Id      uint  `json:"id" binding:"required,gt=0"`
		Enabled *bool `json:"enabled" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := Service.SetState(ctx.Request.Context(), body.Id, *body.Enabled, false); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx)
}
