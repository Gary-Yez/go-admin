package sys_menu

import (
	request2 "gitee.com/mxcker/go-admin/server/core/types/request"
	"gitee.com/mxcker/go-admin/server/core/types/response"
	"github.com/gin-gonic/gin"
)

type controllerStruct struct{}

func (_ *controllerStruct) Get(ctx *gin.Context) {
	req := new(request2.Req)
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

func (_ *controllerStruct) List(ctx *gin.Context) {
	list, total, err := Service.List()
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.List(ctx, Service.ListToTree(list), total)

}

func (_ *controllerStruct) Create(ctx *gin.Context) {
	menu := new(SysMenu)
	err := ctx.ShouldBindJSON(menu)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	err = Service.Create(menu)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, menu)
}

func (_ *controllerStruct) Delete(ctx *gin.Context) {
	req := new(request2.ReqIds)
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

func (_ *controllerStruct) Edit(ctx *gin.Context) {
	data := new(SysMenu)
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
