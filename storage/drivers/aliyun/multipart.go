package aliyun

import (
	"context"
	"errors"
	"github.com/Gary-Yez/go-admin/storage/internal/objectutil"
	storage "github.com/Gary-Yez/go-admin/storage/internal/types"
	sdk "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
)

func (s *Store) session(u storage.MultipartUpload) sdk.InitiateMultipartUploadResult {
	return sdk.InitiateMultipartUploadResult{Bucket: s.bucket.BucketName, Key: u.Key, UploadID: u.ID}
}
func (s *Store) BeginMultipart(ctx context.Context, key string, o storage.MultipartOptions) (storage.MultipartUpload, error) {
	u, err := objectutil.Begin(ctx, key, o)
	if err != nil {
		return u, err
	}
	out, err := s.bucket.InitiateMultipartUpload(key, sdk.WithContext(ctx), sdk.ContentType(o.ContentType))
	if err != nil {
		return u, err
	}
	u.ID = out.UploadID
	return u, nil
}
func (s *Store) UploadPart(ctx context.Context, u storage.MultipartUpload, n int, r io.Reader) (storage.Part, error) {
	size, err := objectutil.PartSize(ctx, u, n)
	if err != nil {
		return storage.Part{}, err
	}
	body, _, done, err := objectutil.Prepare(ctx, r, storage.PutOptions{Size: size})
	if err != nil {
		return storage.Part{}, err
	}
	defer done()
	out, err := s.bucket.UploadPart(s.session(u), body, size, n, sdk.WithContext(ctx))
	if err != nil {
		return storage.Part{}, err
	}
	return storage.Part{Number: n, Size: size, ETag: out.ETag}, nil
}
func (s *Store) CompleteMultipart(ctx context.Context, u storage.MultipartUpload, parts []storage.Part) (storage.Object, error) {
	ordered, err := objectutil.Parts(ctx, u, parts)
	if err != nil {
		return storage.Object{}, err
	}
	items := make([]sdk.UploadPart, len(ordered))
	for i, p := range ordered {
		items[i] = sdk.UploadPart{PartNumber: p.Number, ETag: p.ETag}
	}
	out, err := s.bucket.CompleteMultipartUpload(s.session(u), items, sdk.WithContext(ctx))
	if err != nil {
		return storage.Object{}, err
	}
	return objectutil.Completed(u, out.ETag), nil
}
func (s *Store) AbortMultipart(ctx context.Context, u storage.MultipartUpload) error {
	if err := objectutil.Upload(ctx, u); err != nil {
		return err
	}
	err := s.bucket.AbortMultipartUpload(s.session(u), sdk.WithContext(ctx))
	var e sdk.ServiceError
	if errors.As(err, &e) && e.Code == "NoSuchUpload" {
		return nil
	}
	return err
}
