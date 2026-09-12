package sys_apis

import (
	"github.com/Gary-Yez/go-admin/request"
	"github.com/Gary-Yez/go-admin/response"
	"github.com/gin-gonic/gin"
)

type controllerStruct struct{}

func (_ *controllerStruct) GetInvalidAPIs(ctx *gin.Context) {
	apis, err := Service.GetInvalidAPIs()
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, apis)
}

func (_ *controllerStruct) GetGroups(ctx *gin.Context) {
	groups, err := Service.GetGroups()
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, groups)
}

func (_ *controllerStruct) List(ctx *gin.Context) {
	req := new(ApiListQuery)
	err := ctx.ShouldBind(req)
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
	data := new(ApiEditBody)
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
