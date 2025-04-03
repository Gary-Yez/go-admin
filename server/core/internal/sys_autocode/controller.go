package sys_autocode

import (
	"gitee.com/mxcker/go-admin/server/core/request"
	"gitee.com/mxcker/go-admin/server/core/response"
	"github.com/gin-gonic/gin"
)

type controllerStruct struct{}

func (_ *controllerStruct) Generate(ctx *gin.Context) {
	data := new(GenerateBody)
	err := ctx.ShouldBindJSON(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	err = Service.SaveHistory(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	err = Service.Generate(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, "success")
}

func (_ *controllerStruct) Preview(ctx *gin.Context) {
	data := new(GenerateBody)
	err := ctx.ShouldBindJSON(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	preview, err := Service.GetTemplates(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, preview)
}

func (_ *controllerStruct) History(ctx *gin.Context) {
	req := new(request.ReqList)
	if err := ctx.ShouldBindQuery(req); err != nil {
		response.Error(ctx, err.Error())
		return
	}
	list, total, err := Service.History(req)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.List(ctx, list, total)
}

func (_ *controllerStruct) GetHistory(ctx *gin.Context) {
	req := new(request.Req)
	err := ctx.ShouldBindQuery(req)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	get, err := Service.Get(req)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, get)
}

func (_ *controllerStruct) DeleteHistory(ctx *gin.Context) {
	req := new(request.ReqIds)
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	err = Service.DeleteByIds(req)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, "success")
}
