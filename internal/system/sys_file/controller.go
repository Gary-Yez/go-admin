package sys_file

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gary-Yez/go-admin/internal/system/sys_storage"
	"github.com/Gary-Yez/go-admin/request"
	"github.com/Gary-Yez/go-admin/response"
	"github.com/gin-gonic/gin"
)

type controllerStruct struct{}

func owner(ctx *gin.Context) (uint, bool) {
	user, err := request.GetAuthUser(ctx)
	if err != nil {
		response.Error(ctx, err, 401)
		return 0, false
	}
	return user.UserId, true
}
func bindId(ctx *gin.Context) (uint, bool) {
	var body IdBody
	if err := ctx.ShouldBind(&body); err != nil {
		response.Error(ctx, "请选择文件")
		return 0, false
	}
	return body.Id, true
}
func (*controllerStruct) List(ctx *gin.Context) {
	userId, ok := owner(ctx)
	if !ok {
		return
	}
	req, err := request.GetReqList(ctx)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	list, total, err := Service.List(ctx.Request.Context(), req, userId)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	options, err := sys_storage.Options(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, err)
		return
	}
	policy, err := readUploadPolicy()
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"list": list, "total": total, "storages": options, "policy": policy})
}

// Options 为业务上传提供配置，不查询文件列表或暴露存储凭据。
func (*controllerStruct) Options(ctx *gin.Context) {
	if _, ok := owner(ctx); !ok {
		return
	}
	options, err := sys_storage.Options(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, err)
		return
	}
	policy, err := readUploadPolicy()
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"storages": options, "policy": policy})
}

func (*controllerStruct) Upload(ctx *gin.Context) {
	userId, ok := owner(ctx)
	if !ok {
		return
	}
	policy, err := readUploadPolicy()
	if err != nil {
		response.Error(ctx, err)
		return
	}
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, policy.OrdinaryLimit+(1<<20))
	err = ctx.Request.ParseMultipartForm(1 << 20)
	if ctx.Request.MultipartForm != nil {
		defer ctx.Request.MultipartForm.RemoveAll()
	}
	if err != nil {
		response.Error(ctx, fmt.Sprintf("读取文件失败，普通上传最大支持 %d MiB", policy.OrdinaryLimit>>20))
		return
	}
	if len(ctx.Request.MultipartForm.File["file"]) != 1 {
		response.Error(ctx, "请通过 file 字段提交一个文件")
		return
	}
	file := ctx.Request.MultipartForm.File["file"][0]
	reader, err := file.Open()
	if err != nil {
		response.Error(ctx, err)
		return
	}
	defer reader.Close()
	storageId := uint64(0)
	if value := ctx.PostForm("storage_id"); value != "" {
		storageId, err = strconv.ParseUint(value, 10, 32)
		if err != nil {
			response.Error(ctx, "存储账号编号无效")
			return
		}
	}
	row, err := Service.Upload(ctx.Request.Context(), userId, uint(storageId), file.Filename, file.Size, reader, policy)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, row)
}
func (*controllerStruct) Begin(ctx *gin.Context) {
	userId, ok := owner(ctx)
	if !ok {
		return
	}
	var body BeginBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.Error(ctx, "请提供有效的文件名、大小和修改时间")
		return
	}
	session, err := Service.Begin(ctx.Request.Context(), userId, &body)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, session)
}
func (*controllerStruct) Session(ctx *gin.Context) {
	userId, ok := owner(ctx)
	if !ok {
		return
	}
	id, ok := bindId(ctx)
	if !ok {
		return
	}
	session, err := Service.Session(ctx.Request.Context(), id, userId)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, session)
}
func (*controllerStruct) Part(ctx *gin.Context) {
	userId, ok := owner(ctx)
	if !ok {
		return
	}
	rawId, err := strconv.ParseUint(ctx.Query("id"), 10, 32)
	if err != nil || rawId == 0 {
		response.Error(ctx, "上传会话编号无效")
		return
	}
	number, err := strconv.Atoi(ctx.Query("number"))
	if err != nil {
		response.Error(ctx, "分片编号无效")
		return
	}
	// 服务层先按会话校验 Content-Length，再读取分片；不使用当前配置限制旧会话。
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, ctx.Request.ContentLength)
	part, err := Service.Part(ctx.Request.Context(), uint(rawId), userId, number, ctx.Request.Body, ctx.Request.ContentLength)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, part)
}
func (*controllerStruct) Complete(ctx *gin.Context) {
	userId, ok := owner(ctx)
	if !ok {
		return
	}
	id, ok := bindId(ctx)
	if !ok {
		return
	}
	row, err := Service.Complete(ctx.Request.Context(), id, userId)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, row)
}
func (*controllerStruct) Abort(ctx *gin.Context) {
	userId, ok := owner(ctx)
	if !ok {
		return
	}
	id, ok := bindId(ctx)
	if !ok {
		return
	}
	if err := Service.Delete(ctx.Request.Context(), id, userId, false); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx)
}
func (*controllerStruct) Delete(ctx *gin.Context) {
	var req struct {
		request.ReqIds
		RecordsOnly bool `json:"records_only"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := req.Validate(); err != nil {
		response.Error(ctx, err)
		return
	}
	for i, id := range req.Ids {
		if err := Service.Delete(ctx.Request.Context(), id, 0, req.RecordsOnly); err != nil {
			response.Error(ctx, fmt.Errorf("已删除 %d 个文件，编号 %d 删除失败：%w", i, id, err))
			return
		}
	}
	response.Success(ctx)
}
func (*controllerStruct) Cleanup(ctx *gin.Context) {
	count, err := Service.Cleanup(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, gin.H{"count": count})
}
