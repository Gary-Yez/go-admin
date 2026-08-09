package sys_devtools

import (
	"gitee.com/mxcker/go-admin/request"
	"gitee.com/mxcker/go-admin/response"
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
	req, err := request.GetReqList(ctx)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	list, total, err := Service.History(req)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.List(ctx, list, total)
}

func (_ *controllerStruct) DeleteHistory(ctx *gin.Context) {
	req, err := request.GetReqIds(ctx)
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
