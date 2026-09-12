package sys_devtools

import (
	"github.com/Gary-Yez/go-admin/request"
	"github.com/Gary-Yez/go-admin/response"
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
	req := &HistoryQuery{Page: 1, Limit: 10}
	err := ctx.ShouldBindQuery(req)
	if err != nil {
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
	req, err := request.GetReq(ctx)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	history, err := Service.GetHistory(req)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, history)
}

func (_ *controllerStruct) DeleteHistory(ctx *gin.Context) {
	req := new(DeleteHistoryBody)
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	err = Service.DeleteHistory(req)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, "success")
}

func (_ *controllerStruct) PreviewDeleteHistory(ctx *gin.Context) {
	req := new(DeleteHistoryBody)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.Error(ctx, err.Error())
		return
	}
	plan, err := Service.PreviewDeleteHistory(req.Ids)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, plan)
}
