package sys_admin

import (
	"errors"
	"gitee.com/mxcker/go-admin/server/types"
	request2 "gitee.com/mxcker/go-admin/server/types/request"
	"gitee.com/mxcker/go-admin/server/types/response"
	"github.com/gin-gonic/gin"
	"slices"
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
	req := new(request2.ReqList)
	if err := ctx.ShouldBindQuery(req); err != nil {
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
	data := new(SysAdmin)
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
	req := new(request2.ReqIds)
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	authUser, ok := ctx.MustGet("AuthUser").(types.AuthUser)
	if !ok {
		response.Error(ctx, "登录失效")
		return
	}
	if slices.Contains(req.Ids, authUser.UserId) {
		response.Error(ctx, errors.New("不可以自己删除自己"))
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
	data := new(SysAdmin)
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
