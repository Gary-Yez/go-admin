// Package local 实现受限目录中的本地文件存储。
package local

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/Gary-Yez/go-admin/storage/internal/objectutil"
	storage "github.com/Gary-Yez/go-admin/storage/internal/types"
)

type Config struct {
	Root      string
	PublicURL string
}

// 固定数量的锁供同一进程的新实例共享，避免按目录无限保留锁。
var multipartLocks [64]sync.Mutex

func directoryLock(name string) *sync.Mutex {
	if runtime.GOOS == "windows" {
		name = strings.ToLower(name)
	}
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(name))
	return &multipartLocks[hash.Sum32()%uint32(len(multipartLocks))]
}

type Store struct {
	multipartMu *sync.Mutex
	root        *os.Root
	staging     *os.Root
	publicURL   string
}

var _ storage.Store = (*Store)(nil)
var _ storage.URLProvider = (*Store)(nil)

func New(config Config) (*Store, error) {
	if config.Root == "" {
		return nil, fmt.Errorf("本地存储根目录不能为空")
	}
	if err := objectutil.PublicBase(config.PublicURL); err != nil {
		return nil, err
	}
	absolute, err := filepath.Abs(config.Root)
	if err != nil {
		return nil, err
	}
	if err = os.MkdirAll(absolute, 0750); err != nil {
		return nil, err
	}
	root, err := os.OpenRoot(absolute)
	if err != nil {
		return nil, err
	}
	stagingPath := absolute + ".multipart"
	if err = os.MkdirAll(stagingPath, 0700); err != nil {
		root.Close()
		return nil, err
	}
	staging, err := os.OpenRoot(stagingPath)
	if err != nil {
		root.Close()
		return nil, err
	}
	canonical, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		root.Close()
		staging.Close()
		return nil, err
	}
	return &Store{multipartMu: directoryLock(canonical), root: root, staging: staging, publicURL: config.PublicURL}, nil
}
func (s *Store) Close() error { return errors.Join(s.root.Close(), s.staging.Close()) }
func (s *Store) Put(ctx context.Context, key string, r io.Reader, opts storage.PutOptions) (storage.Object, error) {
	if err := objectutil.Key(ctx, key); err != nil {
		return storage.Object{}, err
	}
	if r == nil || opts.Size < -1 {
		return storage.Object{}, fmt.Errorf("上传内容或大小无效")
	}
	if err := s.root.MkdirAll(path.Dir(key), 0750); err != nil {
		return storage.Object{}, err
	}
	temp := path.Join(path.Dir(key), ".upload-"+rand.Text())
	file, err := s.root.OpenFile(temp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return storage.Object{}, err
	}
	defer s.root.Remove(temp)
	defer file.Close()
	var source io.Reader = objectutil.Reader{Ctx: ctx, Source: r}
	if opts.Size >= 0 && opts.Size < 1<<63-1 {
		source = io.LimitReader(source, opts.Size+1)
	}
	size, err := io.Copy(file, source)
	if err != nil {
		return storage.Object{}, err
	}
	if opts.Size >= 0 && size != opts.Size {
		return storage.Object{}, fmt.Errorf("实际文件大小与声明不一致")
	}
	if err = ctx.Err(); err != nil {
		return storage.Object{}, err
	}
	if err = file.Sync(); err != nil {
		return storage.Object{}, err
	}
	if err = file.Close(); err != nil {
		return storage.Object{}, err
	}
	if err = s.root.Rename(temp, key); err != nil {
		return storage.Object{}, err
	}
	return s.Stat(ctx, key)
}
func (s *Store) Open(ctx context.Context, key string) (io.ReadCloser, storage.Object, error) {
	if err := objectutil.Key(ctx, key); err != nil {
		return nil, storage.Object{}, err
	}
	file, err := s.root.Open(key)
	if err != nil {
		return nil, storage.Object{}, objectutil.Missing(err, errors.Is(err, os.ErrNotExist))
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, storage.Object{}, err
	}
	if !info.Mode().IsRegular() {
		file.Close()
		return nil, storage.Object{}, fmt.Errorf("目标不是普通文件")
	}
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		file.Close()
		return nil, storage.Object{}, err
	}
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		file.Close()
		return nil, storage.Object{}, err
	}
	obj := storage.Object{Key: key, Size: info.Size(), ModifiedAt: info.ModTime(), ContentType: http.DetectContentType(buf[:n])}
	return objectutil.ReadCloser(ctx, file), obj, nil
}
func (s *Store) Stat(ctx context.Context, key string) (storage.Object, error) {
	r, obj, err := s.Open(ctx, key)
	if err == nil {
		err = r.Close()
	}
	return obj, err
}
func (s *Store) Delete(ctx context.Context, key string) error {
	if err := objectutil.Key(ctx, key); err != nil {
		return err
	}
	info, err := s.root.Lstat(key)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("只允许删除普通文件")
	}
	err = s.root.Remove(key)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
func (s *Store) URL(ctx context.Context, key string, opts storage.URLOptions) (string, error) {
	if err := objectutil.URL(ctx, key, opts); err != nil {
		return "", err
	}
	return objectutil.PublicURL(s.publicURL, key, opts)
}
