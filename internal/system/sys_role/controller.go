package sys_role

import (
	"gitee.com/mxcker/go-admin/request"
	"gitee.com/mxcker/go-admin/response"
	"github.com/gin-gonic/gin"
)

type controllerStruct struct{}

func (_ *controllerStruct) Get(ctx *gin.Context) {
	req, err := request.GetReq(ctx)
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
	response.List(ctx, list, total)

}

func (_ *controllerStruct) Create(ctx *gin.Context) {
	data := new(SysRole)
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
	data := new(SysRole)
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

func (_ *controllerStruct) UpdatePermission(ctx *gin.Context) {
	permission := new(SysRole)
	err := ctx.ShouldBindJSON(permission)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	err = Service.UpdatePermission(permission)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, "success")
}
