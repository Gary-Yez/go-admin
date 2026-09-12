package sys_api_token

import (
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_admin"
	"github.com/Gary-Yez/go-admin/internal/system/sys_role"
	"github.com/Gary-Yez/go-admin/request"
	"github.com/Gary-Yez/go-admin/response"
	"github.com/gin-gonic/gin"
)

type controllerStruct struct{}

func (*controllerStruct) List(ctx *gin.Context) {
	req, err := request.GetReqList(ctx)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	list, total, err := Service.List(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	roles, err := Service.RoleOptions(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"list": list, "total": total, "role_options": roles})
}
func (*controllerStruct) Delete(ctx *gin.Context) {
	req, err := request.GetReqIds(ctx)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	if err := Service.DeleteApiTokens(ctx.Request.Context(), req, 0); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "API 密钥已删除")
}

func (_ *controllerStruct) ListApiTokens(ctx *gin.Context) {
	user, err := request.GetAuthUser(ctx)
	if err != nil {
		response.Error(ctx, err, 401)
		return
	}
	req, err := request.GetReqList(ctx)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	db := req.WithFilter(state.DB().Model(&SysApiToken{}).Where("admin_id = ?", user.UserId), []string{"remark", "role_id"})
	var total int64
	if err := db.Count(&total).Error; err != nil {
		response.Error(ctx, err)
		return
	}
	list := []SysApiToken{}
	if err := req.WithPagination(db).Order("id DESC").Find(&list).Error; err != nil {
		response.Error(ctx, err)
		return
	}
	roles := []sys_admin.RoleOption{}
	if err := state.DB().Model(&sys_role.SysRole{}).Select("id", "name").Where("id IN (?)", state.DB().Model(&sys_admin.SysAdminRole{}).Select("role_id").Where("admin_id = ?", user.UserId)).Order("id").Find(&roles).Error; err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"list": list, "total": total, "role_options": roles})
}

func (_ *controllerStruct) SaveApiToken(ctx *gin.Context) {
	user, err := request.GetAuthUser(ctx)
	if err != nil {
		response.Error(ctx, err, 401)
		return
	}
	var body ApiTokenBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.Error(ctx, err)
		return
	}
	token, err := Service.SaveApiToken(user.UserId, &body)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"token": token})
}

func (_ *controllerStruct) DeleteApiToken(ctx *gin.Context) {
	user, err := request.GetAuthUser(ctx)
	if err != nil {
		response.Error(ctx, err, 401)
		return
	}
	req, err := request.GetReq(ctx)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	if err := Service.DeleteApiTokens(ctx.Request.Context(), &request.ReqIds{Ids: []uint{req.Id}}, user.UserId); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx)
}
