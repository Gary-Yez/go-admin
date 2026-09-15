package objectutil

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode"

	storage "github.com/Gary-Yez/go-admin/storage/internal/types"
)

func Key(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if key == "" || strings.ContainsAny(key, "\\:") || strings.ContainsFunc(key, unicode.IsControl) {
		return fmt.Errorf("无效文件 Key")
	}
	for _, part := range strings.Split(key, "/") {
		if part == "" || part == "." || part == ".." {
			return fmt.Errorf("文件 Key 必须为相对路径且不能包含空段或目录穿越")
		}
	}
	return nil
}
func Endpoint(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("无效存储地址")
	}
	if (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return nil, fmt.Errorf("存储地址必须为不含凭据、查询参数的 HTTP(S) 地址")
	}
	return u, nil
}
func PublicBase(raw string) error {
	if raw == "" {
		return nil
	}
	_, err := Endpoint(raw)
	return err
}
func PublicURL(base, key string, options storage.URLOptions) (string, error) {
	if base == "" || options.Expires != 0 || options.DownloadName != "" {
		return "", storage.ErrUnsupported
	}
	u, err := Endpoint(base)
	if err != nil {
		return "", err
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/" + key
	u.RawPath = ""
	return u.String(), nil
}
func URL(ctx context.Context, key string, options storage.URLOptions) error {
	if err := Key(ctx, key); err != nil {
		return err
	}
	if options.Expires < 0 || options.Expires > 7*24*time.Hour || (options.Expires > 0 && options.Expires < time.Second) {
		return fmt.Errorf("签名有效期应为 1 秒至 7 天，公开地址使用 0")
	}
	if strings.ContainsAny(options.DownloadName, "/\\") || strings.ContainsFunc(options.DownloadName, unicode.IsControl) {
		return fmt.Errorf("无效下载文件名")
	}
	return nil
}
func Disposition(name string) string {
	return mime.FormatMediaType("attachment", map[string]string{"filename": name})
}
func Headers(key string, h http.Header) storage.Object {
	size := int64(0)
	fmt.Sscan(h.Get("Content-Length"), &size)
	modified, _ := http.ParseTime(h.Get("Last-Modified"))
	return storage.Object{Key: key, Size: size, ContentType: h.Get("Content-Type"), ETag: strings.Trim(h.Get("ETag"), "\""), ModifiedAt: modified}
}
func Missing(err error, missing bool) error {
	if err != nil && missing {
		return errors.Join(storage.ErrNotFound, err)
	}
	return err
}

type Reader struct {
	Ctx    context.Context
	Source io.Reader
}

func (r Reader) Read(b []byte) (int, error) {
	if err := r.Ctx.Err(); err != nil {
		return 0, err
	}
	return r.Source.Read(b)
}

type readCloser struct {
	Reader
	io.Closer
}

func ReadCloser(ctx context.Context, r io.ReadCloser) io.ReadCloser {
	return readCloser{Reader{ctx, r}, r}
}

// Prepare 为 SDK 提供可重试、已知长度的流；不会关闭调用方的 reader。
func Prepare(ctx context.Context, reader io.Reader, options storage.PutOptions) (io.ReadSeeker, int64, func(), error) {
	noop := func() {}
	if err := ctx.Err(); err != nil {
		return nil, 0, noop, err
	}
	if reader == nil || options.Size < -1 {
		return nil, 0, noop, fmt.Errorf("上传内容或大小无效")
	}
	if strings.ContainsAny(options.ContentType, "\r\n") {
		return nil, 0, noop, fmt.Errorf("无效 Content-Type")
	}
	if seeker, ok := reader.(io.ReadSeeker); ok {
		pos, err := seeker.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, 0, noop, err
		}
		end, err := seeker.Seek(0, io.SeekEnd)
		if err != nil {
			return nil, 0, noop, err
		}
		if _, err = seeker.Seek(pos, io.SeekStart); err != nil {
			return nil, 0, noop, err
		}
		size := end - pos
		if size < 0 || (options.Size >= 0 && options.Size != size) {
			return nil, 0, noop, fmt.Errorf("实际文件大小与声明不一致")
		}
		return &offsetSeeker{seeker: seeker, base: pos}, size, noop, nil
	}
	file, err := os.CreateTemp("", "go-admin-upload-*")
	if err != nil {
		return nil, 0, noop, err
	}
	cleanup := func() { file.Close(); os.Remove(file.Name()) }
	var source io.Reader = Reader{ctx, reader}
	if options.Size >= 0 && options.Size < 1<<63-1 {
		source = io.LimitReader(source, options.Size+1)
	}
	size, err := io.Copy(file, source)
	if err == nil {
		err = ctx.Err()
	}
	if err == nil && options.Size >= 0 && options.Size != size {
		err = fmt.Errorf("实际文件大小与声明不一致")
	}
	if err == nil {
		_, err = file.Seek(0, io.SeekStart)
	}
	if err != nil {
		cleanup()
		return nil, 0, noop, err
	}
	return &offsetSeeker{seeker: file}, size, cleanup, nil
}

type offsetSeeker struct {
	seeker io.ReadSeeker
	base   int64
}

func (s *offsetSeeker) Read(p []byte) (int, error) { return s.seeker.Read(p) }
func (s *offsetSeeker) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekStart {
		offset += s.base
	}
	pos, err := s.seeker.Seek(offset, whence)
	return pos - s.base, err
}
