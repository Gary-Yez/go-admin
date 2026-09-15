package objectutil

import (
	"context"
	"fmt"
	storage "github.com/Gary-Yez/go-admin/storage/internal/types"
	"sort"
	"strings"
)

const MaxParts = 10000

func Begin(ctx context.Context, key string, o storage.MultipartOptions) (storage.MultipartUpload, error) {
	if err := Key(ctx, key); err != nil {
		return storage.MultipartUpload{}, err
	}
	if o.PartSize == 0 {
		o.PartSize = 50 << 20
	}
	if o.Size <= 0 || o.PartSize < 5<<20 || o.PartSize > 1<<30 || (o.Size-1)/o.PartSize+1 > MaxParts || strings.ContainsAny(o.ContentType, "\r\n") {
		return storage.MultipartUpload{}, fmt.Errorf("文件大小或分片配置无效：分片为 5 MiB 至 1 GiB，最多 10000 片")
	}
	return storage.MultipartUpload{Key: key, Size: o.Size, PartSize: o.PartSize, ContentType: o.ContentType}, nil
}
func Upload(ctx context.Context, u storage.MultipartUpload) error {
	if u.ID == "" {
		return fmt.Errorf("上传 ID 不能为空")
	}
	_, err := Begin(ctx, u.Key, storage.MultipartOptions{Size: u.Size, PartSize: u.PartSize, ContentType: u.ContentType})
	if u.PartSize == 0 {
		return fmt.Errorf("上传分片大小不能为空")
	}
	return err
}
func PartSize(ctx context.Context, u storage.MultipartUpload, n int) (int64, error) {
	if err := Upload(ctx, u); err != nil {
		return 0, err
	}
	count := (u.Size-1)/u.PartSize + 1
	if n < 1 || int64(n) > count {
		return 0, fmt.Errorf("无效分片编号")
	}
	size := u.PartSize
	if int64(n) == count {
		size = u.Size - int64(n-1)*u.PartSize
	}
	return size, nil
}
func Parts(ctx context.Context, u storage.MultipartUpload, parts []storage.Part) ([]storage.Part, error) {
	if err := Upload(ctx, u); err != nil {
		return nil, err
	}
	if int64(len(parts)) != (u.Size-1)/u.PartSize+1 {
		return nil, fmt.Errorf("分片尚未全部上传")
	}
	result := append([]storage.Part(nil), parts...)
	sort.Slice(result, func(i, j int) bool { return result[i].Number < result[j].Number })
	for i, p := range result {
		size, err := PartSize(ctx, u, i+1)
		if err != nil {
			return nil, err
		}
		if p.Number != i+1 || p.Size != size || p.ETag == "" {
			return nil, fmt.Errorf("分片回执不完整、重复或大小不正确")
		}
	}
	return result, nil
}
func Completed(u storage.MultipartUpload, etag string) storage.Object {
	return storage.Object{Key: u.Key, Size: u.Size, ContentType: u.ContentType, ETag: etag}
}
