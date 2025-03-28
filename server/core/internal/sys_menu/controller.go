package sys_menu

import (
	"gitee.com/mxcker/go-admin/server/core/request"
	"gitee.com/mxcker/go-admin/server/core/response"
	"github.com/gin-gonic/gin"
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
	list, total, err := service.List()
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.List(ctx, service.ListToTree(list), total)

}

func (_ *Controller) Create(ctx *gin.Context) {
	menu := new(SysMenu)
	err := ctx.ShouldBindJSON(menu)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	err = service.Create(menu)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, menu)
}

func (_ *Controller) Delete(ctx *gin.Context) {
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

func (_ *Controller) Edit(ctx *gin.Context) {
	data := new(SysMenu)
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
