package tencent

import (
	"context"
	"errors"
	"github.com/Gary-Yez/go-admin/storage/internal/objectutil"
	storage "github.com/Gary-Yez/go-admin/storage/internal/types"
	sdk "github.com/tencentyun/cos-go-sdk-v5"
	"io"
)

func (s *Store) BeginMultipart(ctx context.Context, key string, o storage.MultipartOptions) (storage.MultipartUpload, error) {
	u, err := objectutil.Begin(ctx, key, o)
	if err != nil {
		return u, err
	}
	out, _, err := s.client.Object.InitiateMultipartUpload(ctx, key, &sdk.InitiateMultipartUploadOptions{ObjectPutHeaderOptions: &sdk.ObjectPutHeaderOptions{ContentType: o.ContentType}})
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
	out, err := s.client.Object.UploadPart(ctx, u.Key, u.ID, n, body, &sdk.ObjectUploadPartOptions{ContentLength: size})
	if err != nil {
		return storage.Part{}, err
	}
	return storage.Part{Number: n, Size: size, ETag: out.Header.Get("ETag")}, nil
}
func (s *Store) CompleteMultipart(ctx context.Context, u storage.MultipartUpload, parts []storage.Part) (storage.Object, error) {
	ordered, err := objectutil.Parts(ctx, u, parts)
	if err != nil {
		return storage.Object{}, err
	}
	items := make([]sdk.Object, len(ordered))
	for i, p := range ordered {
		items[i] = sdk.Object{PartNumber: p.Number, ETag: p.ETag}
	}
	out, _, err := s.client.Object.CompleteMultipartUpload(ctx, u.Key, u.ID, &sdk.CompleteMultipartUploadOptions{Parts: items})
	if err != nil {
		return storage.Object{}, err
	}
	return objectutil.Completed(u, out.ETag), nil
}
func (s *Store) AbortMultipart(ctx context.Context, u storage.MultipartUpload) error {
	if err := objectutil.Upload(ctx, u); err != nil {
		return err
	}
	_, err := s.client.Object.AbortMultipartUpload(ctx, u.Key, u.ID)
	var e *sdk.ErrorResponse
	if errors.As(err, &e) && e.Code == "NoSuchUpload" {
		return nil
	}
	return err
}
