package sys_file

import (
	"context"
	"errors"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/response"
	"github.com/Gary-Yez/go-admin/storage"
	"github.com/gin-gonic/gin"
)

func textPreviewable(contentType string) bool {
	switch contentType {
	case "application/json", "application/xml", "application/javascript", "application/x-ndjson":
		return true
	}
	return strings.HasPrefix(contentType, "text/")
}

func previewable(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/avif",
		"video/mp4", "video/webm", "video/ogg",
		"audio/mpeg", "audio/mp4", "audio/aac", "audio/wav", "audio/x-wav", "audio/wave", "audio/ogg", "audio/flac", "audio/webm", "application/ogg",
		"application/pdf":
		return true
	}
	return textPreviewable(contentType)
}

func (*controllerStruct) Link(ctx *gin.Context) {
	var body struct {
		Id      uint `json:"id" binding:"required,gt=0"`
		Preview bool `json:"preview"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.Error(ctx, "请选择文件")
		return
	}
	var row SysFile
	if err := state.DB().WithContext(ctx.Request.Context()).Where("status = ?", ready).First(&row, body.Id).Error; err != nil {
		response.Error(ctx, err)
		return
	}
	if body.Preview && !previewable(row.ContentType) {
		response.Error(ctx, "该类型不支持在线预览，请下载查看")
		return
	}
	url, expires, err := fileLink(ctx.Request.Context(), &row, body.Preview)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	ctx.Header("Cache-Control", "no-store")
	response.Success(ctx, gin.H{"url": url, "expires_in": expires})
}

// fileLink 优先直连公开地址，其次使用私有签名；仅本地无公开地址时由 Go 返回内容。
func fileLink(ctx context.Context, row *SysFile, preview bool) (string, int, error) {
	store, err := row.openStore()
	if err != nil {
		return "", 0, err
	}
	defer closeStore(store)
	if provider, ok := store.(storage.URLProvider); ok {
		url, err := provider.URL(ctx, row.ObjectKey, storage.URLOptions{})
		if err == nil {
			return url, 0, nil
		}
		if !errors.Is(err, storage.ErrUnsupported) {
			return "", 0, err
		}
		if row.Engine != string(storage.Local) {
			ttl, err := readLinkTTL()
			if err != nil {
				return "", 0, err
			}
			options := storage.URLOptions{Expires: ttl}
			if !preview {
				options.DownloadName = row.Name
			}
			url, err := provider.URL(ctx, row.ObjectKey, options)
			return url, int(ttl / time.Second), err
		}
	} else if row.Engine != string(storage.Local) {
		return "", 0, storage.ErrUnsupported
	}
	return localFileLink(ctx, row.Id, preview)
}

func (*controllerStruct) Content(ctx *gin.Context) {
	id, err := resolveLocalLink(ctx.Param("ticket"))
	if err != nil {
		ctx.String(http.StatusNotFound, "下载链接无效或已过期")
		return
	}
	var row SysFile
	if err := state.DB().WithContext(ctx.Request.Context()).Where("status = ?", ready).First(&row, id).Error; err != nil {
		ctx.String(http.StatusNotFound, "文件不存在")
		return
	}
	store, err := row.openStore()
	if err != nil {
		ctx.String(http.StatusServiceUnavailable, "文件存储暂时不可用")
		return
	}
	defer closeStore(store)
	reader, object, err := store.Open(ctx.Request.Context(), row.ObjectKey)
	if err != nil {
		ctx.String(http.StatusNotFound, "文件读取失败")
		return
	}
	defer reader.Close()
	disposition := "attachment"
	contentType := "application/octet-stream"
	if ctx.Query("download") != "1" && previewable(row.ContentType) {
		disposition = "inline"
		contentType = row.ContentType
		if textPreviewable(contentType) {
			// HTML、XML、脚本等统一作为纯文本返回，不作为页面执行。
			contentType = "text/plain; charset=utf-8"
		}
	}
	ctx.Header("Content-Disposition", mime.FormatMediaType(disposition, map[string]string{"filename": row.Name}))
	ctx.Header("X-Content-Type-Options", "nosniff")
	ctx.Header("Cache-Control", "private, no-store")
	ctx.Header("Referrer-Policy", "no-referrer")
	ctx.Header("Content-Security-Policy", "default-src 'none'; sandbox")
	// Store.Open 是流式接口，不把整个文件读入内存。可寻址的本地文件同时支持 Range。
	if seeker, ok := reader.(io.ReadSeeker); ok {
		ctx.Header("Content-Type", contentType)
		http.ServeContent(ctx.Writer, ctx.Request, strings.TrimSpace(row.Name), object.ModifiedAt, seeker)
		return
	}
	ctx.DataFromReader(http.StatusOK, object.Size, contentType, reader, nil)
}

// PreviewLink 为已取得业务授权的调用方生成图片地址；不提供文件归属判断。
func (f *SysFile) PreviewLink(ctx context.Context) (string, error) {
	if f.Status != ready || !strings.HasPrefix(f.ContentType, "image/") || !previewable(f.ContentType) {
		return "", errors.New("文件未完成上传或不支持图片预览")
	}
	url, _, err := fileLink(ctx, f, true)
	return url, err
}

// ReadContent 在回调结束时关闭读取流和存储资源，供业务校验真实文件内容。
func (f *SysFile) ReadContent(ctx context.Context, read func(io.Reader) error) error {
	store, err := f.openStore()
	if err != nil {
		return err
	}
	defer closeStore(store)
	reader, _, err := store.Open(ctx, f.ObjectKey)
	if err != nil {
		return err
	}
	defer reader.Close()
	return read(reader)
}
