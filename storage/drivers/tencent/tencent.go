// Package tencent 实现腾讯云 COS 存储。
package tencent

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Gary-Yez/go-admin/storage/internal/objectutil"
	storage "github.com/Gary-Yez/go-admin/storage/internal/types"
	sdk "github.com/tencentyun/cos-go-sdk-v5"
)

type Config struct{ BucketURL, SecretID, SecretKey, PublicURL string }
type Store struct {
	client    *sdk.Client
	publicURL string
}

var _ storage.Store = (*Store)(nil)
var _ storage.URLProvider = (*Store)(nil)

func New(c Config) (*Store, error) {
	if c.SecretID == "" || c.SecretKey == "" {
		return nil, fmt.Errorf("COS 凭据不能为空")
	}
	u, err := objectutil.Endpoint(c.BucketURL)
	if err != nil {
		return nil, err
	}
	if err = objectutil.PublicBase(c.PublicURL); err != nil {
		return nil, err
	}
	client := sdk.NewClient(&sdk.BaseURL{BucketURL: u}, &http.Client{Transport: &sdk.AuthorizationTransport{SecretID: c.SecretID, SecretKey: c.SecretKey}})
	return &Store{client, c.PublicURL}, nil
}
func translate(err error) error {
	var e *sdk.ErrorResponse
	if errors.As(err, &e) {
		return objectutil.Missing(err, e.Code == "NoSuchKey" || (e.Code == "" && e.Response != nil && e.Response.StatusCode == 404))
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
	resp, err := s.client.Object.Put(ctx, key, body, &sdk.ObjectPutOptions{ObjectPutHeaderOptions: &sdk.ObjectPutHeaderOptions{ContentType: o.ContentType, ContentLength: size}})
	if err != nil {
		return storage.Object{}, translate(err)
	}
	obj := objectutil.Headers(key, resp.Header)
	obj.Size = size
	obj.ContentType = o.ContentType
	obj.ModifiedAt = time.Now().UTC()
	return obj, nil
}
func (s *Store) Open(ctx context.Context, key string) (io.ReadCloser, storage.Object, error) {
	if err := objectutil.Key(ctx, key); err != nil {
		return nil, storage.Object{}, err
	}
	resp, err := s.client.Object.Get(ctx, key, nil)
	if err != nil {
		return nil, storage.Object{}, translate(err)
	}
	return resp.Body, objectutil.Headers(key, resp.Header), nil
}
func (s *Store) Stat(ctx context.Context, key string) (storage.Object, error) {
	if err := objectutil.Key(ctx, key); err != nil {
		return storage.Object{}, err
	}
	resp, err := s.client.Object.Head(ctx, key, nil)
	if err != nil {
		return storage.Object{}, translate(err)
	}
	return objectutil.Headers(key, resp.Header), nil
}
func (s *Store) Delete(ctx context.Context, key string) error {
	if err := objectutil.Key(ctx, key); err != nil {
		return err
	}
	_, err := s.client.Object.Delete(ctx, key)
	err = translate(err)
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
	opts := &sdk.ObjectGetOptions{}
	if o.DownloadName != "" {
		opts.ResponseContentDisposition = objectutil.Disposition(o.DownloadName)
	}
	u, err := s.client.Object.GetPresignedURL2(ctx, http.MethodGet, key, o.Expires, opts)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
