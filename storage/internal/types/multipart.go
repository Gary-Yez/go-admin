package types

// MultipartUpload 必须由业务服务端保存并与用户、存储名称绑定，不能信任客户端传入的内容。
// Size、PartSize 在开始时确定；Key、ID 和 ContentType 之后不得修改。
type MultipartUpload struct {
	ID          string
	Key         string
	Size        int64
	PartSize    int64
	ContentType string
}
type MultipartOptions struct {
	Size        int64 // 必须提供大于 0 的文件总大小。
	PartSize    int64 // 0 使用 50 MiB；允许 5 MiB 至 1 GiB，最多 10000 片。
	ContentType string
}
type Part struct {
	Number int // 从 1 开始。
	Size   int64
	ETag   string
}
