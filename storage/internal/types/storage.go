// Package types 定义文件存储能力，不依赖数据库、HTTP 框架或用户权限。
package types

import (
	"context"
	"errors"
	"io"
	"time"
)

var (
	// ErrNotFound 表示对象不存在，适配器应支持 errors.Is 判断。
	ErrNotFound = errors.New("文件不存在")
	// ErrUnsupported 表示存储不支持请求的可选能力。
	ErrUnsupported = errors.New("存储不支持此操作")
)

// Object 为文件元数据，Key 是存储内部的相对标识，不是磁盘路径或访问地址。
type Object struct {
	Key         string
	Size        int64
	ContentType string
	ETag        string // 存储返回的版本标识，不保证为文件内容的 MD5。
	ModifiedAt  time.Time
}

type PutOptions struct {
	ContentType string
	Size        int64 // 已知大小；未知时传 -1，0 表示空文件。
}

// Store 由具体存储适配器实现；实现须支持并发调用并响应 context 取消。
// Key 使用斜杠分隔的相对标识，适配器必须拒绝绝对路径及目录穿越。
// 适配器保留原始错误供排查，不处理业务归属、权限或数据库记录。
type Store interface {
	// Put 流式写入对象，同 Key 默认覆盖。调用方应使用唯一 Key 防止误覆盖。
	// 不关闭调用方传入的 reader，失败时由适配器负责清理未完成上传。
	Put(ctx context.Context, key string, reader io.Reader, options PutOptions) (Object, error)
	// Open 返回读取流和元数据；调用方必须关闭返回的流。
	Open(ctx context.Context, key string) (io.ReadCloser, Object, error)
	Stat(ctx context.Context, key string) (Object, error)
	// Delete 应当幂等，对象已不存在时也返回 nil。
	Delete(ctx context.Context, key string) error

	// 分片会话和回执由调用方保存；同编号重传覆盖旧分片。
	BeginMultipart(context.Context, string, MultipartOptions) (MultipartUpload, error)
	UploadPart(context.Context, MultipartUpload, int, io.Reader) (Part, error)
	// 合并/取消不得与分片写入并发。超时后先核实结果，再决定重试。
	CompleteMultipart(context.Context, MultipartUpload, []Part) (Object, error)
	AbortMultipart(context.Context, MultipartUpload) error
}

type URLOptions struct {
	// Expires > 0 请求短期签名地址；0 请求公开地址；负值应报错。
	// 公开地址需要存储本身已配置公开访问，不能通过此选项绕过权限。
	Expires time.Duration
	// DownloadName 非空时请求下载文件名，适配器负责正确编码。
	DownloadName string
}

// URLProvider 是可选能力。本地存储未配置 HTTP 访问时无需实现。
// 获取地址不代表取得业务授权，调用方必须事先完成权限检查。
// 不支持相应 URL 选项时应返回 ErrUnsupported，不能悄悄降级为公开地址。
type URLProvider interface {
	URL(ctx context.Context, key string, options URLOptions) (string, error)
}
