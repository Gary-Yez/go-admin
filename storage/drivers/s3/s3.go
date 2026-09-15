// Package s3 实现 AWS S3，以及使用 S3 协议的存储。
package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Gary-Yez/go-admin/storage/internal/objectutil"
	storage "github.com/Gary-Yez/go-admin/storage/internal/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	sdk "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

type Config struct {
	Region      string
	Bucket      string
	Endpoint    string // AWS 留空；兼容服务使用完整 HTTP(S) 地址。
	AccessKey   string
	SecretKey   string
	Credentials aws.CredentialsProvider // 可选：IAM 等外部凭据提供者，优先于静态凭据。
	PathStyle   bool
	PublicURL   string // 可选：已配置公开读取的桶域名或 CDN 前缀。
}
type Store struct {
	client            *sdk.Client
	bucket, publicURL string
}

var _ storage.Store = (*Store)(nil)
var _ storage.URLProvider = (*Store)(nil)

func New(c Config) (*Store, error) {
	if c.Bucket == "" || c.Region == "" {
		return nil, fmt.Errorf("S3 Bucket 和 Region 不能为空")
	}
	if err := objectutil.PublicBase(c.PublicURL); err != nil {
		return nil, err
	}
	if c.Endpoint != "" {
		if _, err := objectutil.Endpoint(c.Endpoint); err != nil {
			return nil, err
		}
	}
	provider := c.Credentials
	if provider == nil {
		if c.AccessKey == "" || c.SecretKey == "" {
			return nil, fmt.Errorf("S3 凭据不能为空")
		}
		provider = credentials.NewStaticCredentialsProvider(c.AccessKey, c.SecretKey, "")
	}
	client := sdk.NewFromConfig(aws.Config{
		HTTPClient:  http.DefaultClient,
		Region:      c.Region,
		Credentials: aws.NewCredentialsCache(provider),
	}, func(o *sdk.Options) {
		o.UsePathStyle = c.PathStyle
		o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
		o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
		if c.Endpoint != "" {
			o.BaseEndpoint = aws.String(c.Endpoint)
		}
	})
	return &Store{client, c.Bucket, c.PublicURL}, nil
}
func translate(err error) error {
	var api smithy.APIError
	if errors.As(err, &api) {
		return objectutil.Missing(err, api.ErrorCode() == "NoSuchKey" || api.ErrorCode() == "NotFound")
	}
	return err
}
func (s *Store) Put(ctx context.Context, key string, r io.Reader, opts storage.PutOptions) (storage.Object, error) {
	if err := objectutil.Key(ctx, key); err != nil {
		return storage.Object{}, err
	}
	body, size, cleanup, err := objectutil.Prepare(ctx, r, opts)
	if err != nil {
		return storage.Object{}, err
	}
	defer cleanup()
	out, err := s.client.PutObject(ctx, &sdk.PutObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key), Body: body, ContentLength: aws.Int64(size), ContentType: aws.String(opts.ContentType)})
	if err != nil {
		return storage.Object{}, translate(err)
	}
	return storage.Object{Key: key, Size: size, ContentType: opts.ContentType, ETag: aws.ToString(out.ETag), ModifiedAt: time.Now().UTC()}, nil
}
func (s *Store) Open(ctx context.Context, key string) (io.ReadCloser, storage.Object, error) {
	if err := objectutil.Key(ctx, key); err != nil {
		return nil, storage.Object{}, err
	}
	out, err := s.client.GetObject(ctx, &sdk.GetObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key)})
	if err != nil {
		return nil, storage.Object{}, translate(err)
	}
	return out.Body, storage.Object{Key: key, Size: aws.ToInt64(out.ContentLength), ContentType: aws.ToString(out.ContentType), ETag: aws.ToString(out.ETag), ModifiedAt: aws.ToTime(out.LastModified)}, nil
}
func (s *Store) Stat(ctx context.Context, key string) (storage.Object, error) {
	if err := objectutil.Key(ctx, key); err != nil {
		return storage.Object{}, err
	}
	out, err := s.client.HeadObject(ctx, &sdk.HeadObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key)})
	if err != nil {
		return storage.Object{}, translate(err)
	}
	return storage.Object{Key: key, Size: aws.ToInt64(out.ContentLength), ContentType: aws.ToString(out.ContentType), ETag: aws.ToString(out.ETag), ModifiedAt: aws.ToTime(out.LastModified)}, nil
}
func (s *Store) Delete(ctx context.Context, key string) error {
	if err := objectutil.Key(ctx, key); err != nil {
		return err
	}
	_, err := s.client.DeleteObject(ctx, &sdk.DeleteObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key)})
	err = translate(err)
	if errors.Is(err, storage.ErrNotFound) {
		return nil
	}
	return err
}
func (s *Store) URL(ctx context.Context, key string, opts storage.URLOptions) (string, error) {
	if err := objectutil.URL(ctx, key, opts); err != nil {
		return "", err
	}
	if opts.Expires == 0 {
		return objectutil.PublicURL(s.publicURL, key, opts)
	}
	input := &sdk.GetObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key)}
	if opts.DownloadName != "" {
		input.ResponseContentDisposition = aws.String(objectutil.Disposition(opts.DownloadName))
	}
	out, err := sdk.NewPresignClient(s.client).PresignGetObject(ctx, input, func(o *sdk.PresignOptions) { o.Expires = opts.Expires })
	if err != nil {
		return "", err
	}
	return out.URL, nil
}
