package sys_admin

import (
	"errors"
	request2 "github.com/Gary-Yez/go-admin/request"
	"github.com/Gary-Yez/go-admin/response"
	"github.com/gin-gonic/gin"
	"slices"
)

type controllerStruct struct{}

func (_ *controllerStruct) List(ctx *gin.Context) {
	req := new(AdminListQuery)
	err := ctx.ShouldBind(req)
	if err == nil {
		err = req.Validate()
	}
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	list, total, err := Service.List(req)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	options, err := Service.RoleOptions()
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, gin.H{
		"list":         list,
		"total":        total,
		"role_options": options,
	})
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
	req, err := request2.GetReqIds(ctx)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	authUser, err := request2.GetAuthUser(ctx)
	if err != nil {
		response.Error(ctx, err.Error(), 401)
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
	authUser, err := request2.GetAuthUser(ctx)
	if err != nil {
		response.Error(ctx, err.Error(), 401)
		return
	}
	if data.Id == authUser.UserId && data.Status != 1 {
		response.Error(ctx, "不能禁用当前登录账号")
		return
	}
	err = Service.Update(data)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}
