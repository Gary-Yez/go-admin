package s3

import (
	"context"
	"errors"
	"github.com/Gary-Yez/go-admin/storage/internal/objectutil"
	storage "github.com/Gary-Yez/go-admin/storage/internal/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	sdk "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"io"
)

func (s *Store) BeginMultipart(ctx context.Context, key string, o storage.MultipartOptions) (storage.MultipartUpload, error) {
	u, err := objectutil.Begin(ctx, key, o)
	if err != nil {
		return u, err
	}
	out, err := s.client.CreateMultipartUpload(ctx, &sdk.CreateMultipartUploadInput{Bucket: aws.String(s.bucket), Key: aws.String(key), ContentType: aws.String(o.ContentType)})
	if err != nil {
		return u, err
	}
	u.ID = aws.ToString(out.UploadId)
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
	out, err := s.client.UploadPart(ctx, &sdk.UploadPartInput{Bucket: aws.String(s.bucket), Key: aws.String(u.Key), UploadId: aws.String(u.ID), PartNumber: aws.Int32(int32(n)), Body: body, ContentLength: aws.Int64(size)})
	if err != nil {
		return storage.Part{}, err
	}
	return storage.Part{Number: n, Size: size, ETag: aws.ToString(out.ETag)}, nil
}
func (s *Store) CompleteMultipart(ctx context.Context, u storage.MultipartUpload, parts []storage.Part) (storage.Object, error) {
	ordered, err := objectutil.Parts(ctx, u, parts)
	if err != nil {
		return storage.Object{}, err
	}
	items := make([]types.CompletedPart, len(ordered))
	for i, p := range ordered {
		items[i] = types.CompletedPart{PartNumber: aws.Int32(int32(p.Number)), ETag: aws.String(p.ETag)}
	}
	out, err := s.client.CompleteMultipartUpload(ctx, &sdk.CompleteMultipartUploadInput{Bucket: aws.String(s.bucket), Key: aws.String(u.Key), UploadId: aws.String(u.ID), MultipartUpload: &types.CompletedMultipartUpload{Parts: items}})
	if err != nil {
		return storage.Object{}, err
	}
	return objectutil.Completed(u, aws.ToString(out.ETag)), nil
}
func (s *Store) AbortMultipart(ctx context.Context, u storage.MultipartUpload) error {
	if err := objectutil.Upload(ctx, u); err != nil {
		return err
	}
	_, err := s.client.AbortMultipartUpload(ctx, &sdk.AbortMultipartUploadInput{Bucket: aws.String(s.bucket), Key: aws.String(u.Key), UploadId: aws.String(u.ID)})
	var api smithy.APIError
	if errors.As(err, &api) && api.ErrorCode() == "NoSuchUpload" {
		return nil
	}
	return err
}
