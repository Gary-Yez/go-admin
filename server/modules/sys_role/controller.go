package sys_role

import (
	"gitee.com/mxcker/go-admin/server/models"
	"gitee.com/mxcker/go-admin/server/models/request"
	"gitee.com/mxcker/go-admin/server/modules/sys_menu"
	"gitee.com/mxcker/go-admin/server/utils/response"
	"github.com/gin-gonic/gin"
)

var service = new(Service)
var menuService = new(sys_menu.Service)

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
	//for _, role := range list {
	//	role.Menus = menuService.ListToTree(role.Menus)
	//}
	response.List(ctx, list, total)

}

func (_ *Controller) Create(ctx *gin.Context) {
	data := new(models.SysRole)
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
	err = service.DeleteByIds(req.Ids)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, "success")
}

func (_ *Controller) Edit(ctx *gin.Context) {
	data := new(models.SysRole)
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

func (_ *Controller) UpdatePermission(ctx *gin.Context) {
	permission := new(models.SysRole)

	err := ctx.ShouldBindJSON(permission)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	err = service.UpdatePermission(permission)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, "success")
}
