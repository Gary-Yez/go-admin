// Package aliyun 实现阿里云 OSS 存储。
package aliyun

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Gary-Yez/go-admin/storage/internal/objectutil"
	storage "github.com/Gary-Yez/go-admin/storage/internal/types"
	sdk "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Config struct{ Endpoint, Region, Bucket, AccessKey, SecretKey, PublicURL string }
type Store struct {
	bucket    *sdk.Bucket
	publicURL string
}

var _ storage.Store = (*Store)(nil)
var _ storage.URLProvider = (*Store)(nil)

func New(c Config) (*Store, error) {
	if c.Bucket == "" || c.Region == "" || c.AccessKey == "" || c.SecretKey == "" {
		return nil, fmt.Errorf("OSS 桶、地域及凭据不能为空")
	}
	if _, err := objectutil.Endpoint(c.Endpoint); err != nil {
		return nil, err
	}
	if err := objectutil.PublicBase(c.PublicURL); err != nil {
		return nil, err
	}
	client, err := sdk.New(c.Endpoint, c.AccessKey, c.SecretKey, sdk.HTTPClient(http.DefaultClient), sdk.Region(c.Region), sdk.AuthVersion(sdk.AuthV4))
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(c.Bucket)
	if err != nil {
		return nil, err
	}
	return &Store{bucket, c.PublicURL}, nil
}
func translate(err error) error {
	var e sdk.ServiceError
	if errors.As(err, &e) {
		return objectutil.Missing(err, e.Code == "NoSuchKey" || (e.Code == "" && e.StatusCode == 404))
	}
	return err
}
func (s *Store) Put(ctx context.Context, key string, r io.Reader, o storage.PutOptions) (storage.Object, error) {
	if err := objectutil.Key(ctx, key); err != nil {
		return storage.Object{}, err
	}
	body, size, done, err := objectutil.Prepare(ctx, r, o)
	if err != nil {
		return storage.Object{}, err
	}
	defer done()
	var h http.Header
	err = s.bucket.PutObject(key, body, sdk.WithContext(ctx), sdk.ContentType(o.ContentType), sdk.ContentLength(size), sdk.GetResponseHeader(&h))
	if err != nil {
		return storage.Object{}, translate(err)
	}
	obj := objectutil.Headers(key, h)
	obj.Size = size
	obj.ContentType = o.ContentType
	obj.ModifiedAt = time.Now().UTC()
	return obj, nil
}
func (s *Store) Open(ctx context.Context, key string) (io.ReadCloser, storage.Object, error) {
	if err := objectutil.Key(ctx, key); err != nil {
		return nil, storage.Object{}, err
	}
	var h http.Header
	body, err := s.bucket.GetObject(key, sdk.WithContext(ctx), sdk.GetResponseHeader(&h))
	if err != nil {
		return nil, storage.Object{}, translate(err)
	}
	return body, objectutil.Headers(key, h), nil
}
func (s *Store) Stat(ctx context.Context, key string) (storage.Object, error) {
	if err := objectutil.Key(ctx, key); err != nil {
		return storage.Object{}, err
	}
	h, err := s.bucket.GetObjectDetailedMeta(key, sdk.WithContext(ctx))
	return objectutil.Headers(key, h), translate(err)
}
func (s *Store) Delete(ctx context.Context, key string) error {
	if err := objectutil.Key(ctx, key); err != nil {
		return err
	}
	err := translate(s.bucket.DeleteObject(key, sdk.WithContext(ctx)))
	if errors.Is(err, storage.ErrNotFound) {
		return nil
	}
	return err
}
func (s *Store) URL(ctx context.Context, key string, o storage.URLOptions) (string, error) {
	if err := objectutil.URL(ctx, key, o); err != nil {
		return "", err
	}
	if o.Expires == 0 {
		return objectutil.PublicURL(s.publicURL, key, o)
	}
	opts := []sdk.Option{}
	if o.DownloadName != "" {
		opts = append(opts, sdk.ResponseContentDisposition(objectutil.Disposition(o.DownloadName)))
	}
	return s.bucket.SignURL(key, sdk.HTTPGet, int64(o.Expires/time.Second), opts...)
}
