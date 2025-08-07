package sys_cron_job

import (
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/utils/request"
	"gitee.com/mxcker/go-admin/server/utils/response"
	"github.com/gin-gonic/gin"
)

type controllerStruct struct{}

func (_ *controllerStruct) GetHandlers(ctx *gin.Context) {
	response.Success(ctx, global.Timer.GetHandlers())
}

func (_ *controllerStruct) GetLogs(ctx *gin.Context) {
	req, err := request.GetReqList(ctx)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	list, total, err := Service.GetLogs(req)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.List(ctx, list, total)
}

func (_ *controllerStruct) Get(ctx *gin.Context) {
	req, err := request.GetReq(ctx)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	get, err := Service.Get(req)
	if err != nil {
		return
	}
	response.Success(ctx, get)
}

func (_ *controllerStruct) List(ctx *gin.Context) {
	req, err := request.GetReqList(ctx)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	list, total, err := Service.List(req)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.List(ctx, list, total)

}

func (_ *controllerStruct) Create(ctx *gin.Context) {
	data := new(SysCronJob)
	err := ctx.ShouldBindJSON(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	err = Service.Create(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}

func (_ *controllerStruct) Delete(ctx *gin.Context) {
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

func (_ *controllerStruct) Edit(ctx *gin.Context) {
	data := new(SysCronJob)
	err := ctx.ShouldBindJSON(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	err = Service.Update(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}
