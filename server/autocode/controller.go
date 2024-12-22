package autocode

import (
	"gitee.com/mxcker/go-admin/server/models/request"
	"gitee.com/mxcker/go-admin/server/utils/response"
	"github.com/gin-gonic/gin"
)

var service = new(Service)

type Controller struct{}

func (_ *Controller) Generate(ctx *gin.Context) {
	data := new(GenerateBody)
	err := ctx.ShouldBindJSON(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	err = service.SaveHistory(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	err = service.Generate(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, "success")
}

func (_ *Controller) Preview(ctx *gin.Context) {
	data := new(GenerateBody)
	err := ctx.ShouldBindJSON(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	preview, err := service.GetTemplates(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, preview)
}

func (_ *Controller) History(ctx *gin.Context) {
	req := new(request.ReqList)
	if err := ctx.ShouldBindQuery(req); err != nil {
		response.Error(ctx, err.Error())
		return
	}
	list, total, err := service.History(req)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.List(ctx, list, total)
}

func (_ *Controller) GetHistory(ctx *gin.Context) {
	req := new(request.Req)
	err := ctx.ShouldBindQuery(req)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	get, err := service.Get(req.Id)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, get)
}

func (_ *Controller) DeleteHistory(ctx *gin.Context) {
	req := new(request.ReqIds)
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	err = service.DeleteByIds(req.Ids)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, "success")
}
