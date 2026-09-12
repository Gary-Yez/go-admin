package sys_menu

import (
	"github.com/Gary-Yez/go-admin/request"
	"github.com/Gary-Yez/go-admin/response"
	"github.com/gin-gonic/gin"
)

type controllerStruct struct{}

func (_ *controllerStruct) Sort(ctx *gin.Context) {
	req := new(SortBody)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.Error(ctx, err.Error())
		return
	}
	if err := Service.Sort(req); err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, "排序已保存")
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
