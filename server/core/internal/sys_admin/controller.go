package sys_admin

import (
	"gitee.com/mxcker/go-admin/server/core/common"
	"gitee.com/mxcker/go-admin/server/core/request"
	"gitee.com/mxcker/go-admin/server/core/response"
	"github.com/gin-gonic/gin"
	"slices"
)

var controller = new(Controller)

type Controller struct{}

func (_ *Controller) Get(ctx *gin.Context) {
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

func (_ *Controller) List(ctx *gin.Context) {
	req := new(request.ReqList)
	if err := ctx.ShouldBindQuery(req); err != nil {
		response.Error(ctx, err.Error())
		return
	}
	list, total, err := service.List(req)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.List(ctx, list, total)

}

func (_ *Controller) Create(ctx *gin.Context) {
	data := new(SysAdmin)
	err := ctx.ShouldBindJSON(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	err = service.Create(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}

func (_ *Controller) Delete(ctx *gin.Context) {
	req := new(request.ReqIds)
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	authUser, ok := ctx.MustGet("AuthUser").(common.AuthUser)
	if !ok {
		response.Error(ctx, "登录失效")
		return
	}
	if slices.Contains(req.Ids, authUser.UserId) {
		//response.Error(ctx, errors.New("不可以自己删除自己"))
		//return
	}
	err = service.DeleteByIds(req.Ids)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, "success")
}

func (_ *Controller) Edit(ctx *gin.Context) {
	data := new(SysAdmin)
	err := ctx.ShouldBindJSON(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	err = service.Update(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}
