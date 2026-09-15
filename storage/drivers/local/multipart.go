package local

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/Gary-Yez/go-admin/storage/internal/objectutil"
	storage "github.com/Gary-Yez/go-admin/storage/internal/types"
)

// 分片和会话存放在独立暂存根目录，避免暴露在公开文件目录中。
func sessionDir(u storage.MultipartUpload) (string, error) {
	id, err := hex.DecodeString(u.ID)
	if err != nil || len(id) != 16 || hex.EncodeToString(id) != u.ID {
		return "", fmt.Errorf("无效本地上传 ID")
	}
	return u.ID, nil
}
func (s *Store) load(ctx context.Context, u storage.MultipartUpload) (string, error) {
	if err := objectutil.Upload(ctx, u); err != nil {
		return "", err
	}
	if err := objectutil.Key(ctx, u.Key); err != nil {
		return "", err
	}
	dir, err := sessionDir(u)
	if err != nil {
		return "", err
	}
	f, err := s.staging.Open(dir + "/session.json")
	if err != nil {
		return "", err
	}
	defer f.Close()
	var saved storage.MultipartUpload
	if err = json.NewDecoder(io.LimitReader(f, 16384)).Decode(&saved); err != nil {
		return "", err
	}
	if saved != u {
		return "", fmt.Errorf("上传会话与原始配置不一致")
	}
	return dir, nil
}
func (s *Store) BeginMultipart(ctx context.Context, key string, o storage.MultipartOptions) (storage.MultipartUpload, error) {
	if err := objectutil.Key(ctx, key); err != nil {
		return storage.MultipartUpload{}, err
	}
	u, err := objectutil.Begin(ctx, key, o)
	if err != nil {
		return u, err
	}
	s.multipartMu.Lock()
	defer s.multipartMu.Unlock()
	var id [16]byte
	if _, err = rand.Read(id[:]); err != nil {
		return u, err
	}
	u.ID = hex.EncodeToString(id[:])
	dir, _ := sessionDir(u)
	if err = s.staging.Mkdir(dir, 0700); err != nil {
		return u, err
	}
	f, err := s.staging.OpenFile(dir+"/session.json", os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err == nil {
		err = json.NewEncoder(f).Encode(u)
		if err == nil {
			err = f.Sync()
		}
		err = errors.Join(err, f.Close())
	}
	if err != nil {
		s.staging.Remove(dir + "/session.json")
		s.staging.Remove(dir)
		return storage.MultipartUpload{}, err
	}
	return u, nil
}
func (s *Store) UploadPart(ctx context.Context, u storage.MultipartUpload, n int, r io.Reader) (storage.Part, error) {
	size, err := objectutil.PartSize(ctx, u, n)
	if err != nil {
		return storage.Part{}, err
	}
	if r == nil {
		return storage.Part{}, fmt.Errorf("分片内容不能为空")
	}
	s.multipartMu.Lock()
	defer s.multipartMu.Unlock()
	dir, err := s.load(ctx, u)
	if err != nil {
		return storage.Part{}, err
	}
	temp := dir + "/pending-" + rand.Text()
	f, err := s.staging.OpenFile(temp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return storage.Part{}, err
	}
	defer s.staging.Remove(temp)
	defer f.Close()
	hash := sha256.New()
	written, err := io.Copy(io.MultiWriter(f, hash), io.LimitReader(objectutil.Reader{Ctx: ctx, Source: r}, size+1))
	if err != nil {
		return storage.Part{}, err
	}
	if written != size {
		return storage.Part{}, fmt.Errorf("实际分片大小不正确")
	}
	if err = ctx.Err(); err != nil {
		return storage.Part{}, err
	}
	if err = f.Sync(); err != nil {
		return storage.Part{}, err
	}
	if err = f.Close(); err != nil {
		return storage.Part{}, err
	}
	if err = s.staging.Rename(temp, dir+"/"+strconv.Itoa(n)); err != nil {
		return storage.Part{}, err
	}
	return storage.Part{Number: n, Size: size, ETag: hex.EncodeToString(hash.Sum(nil))}, nil
}
func (s *Store) CompleteMultipart(ctx context.Context, u storage.MultipartUpload, parts []storage.Part) (storage.Object, error) {
	ordered, err := objectutil.Parts(ctx, u, parts)
	if err != nil {
		return storage.Object{}, err
	}
	s.multipartMu.Lock()
	defer s.multipartMu.Unlock()
	dir, err := s.load(ctx, u)
	if err != nil {
		return storage.Object{}, err
	}
	if err = s.root.MkdirAll(path.Dir(u.Key), 0750); err != nil {
		return storage.Object{}, err
	}
	temp := path.Join(path.Dir(u.Key), ".upload-"+rand.Text())
	f, err := s.root.OpenFile(temp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return storage.Object{}, err
	}
	defer s.root.Remove(temp)
	defer f.Close()
	for _, p := range ordered {
		part, err := s.staging.Open(dir + "/" + strconv.Itoa(p.Number))
		if err != nil {
			return storage.Object{}, err
		}
		hash := sha256.New()
		size, copyErr := io.Copy(io.MultiWriter(f, hash), io.LimitReader(objectutil.Reader{Ctx: ctx, Source: part}, p.Size+1))
		closeErr := part.Close()
		if err = errors.Join(copyErr, closeErr); err != nil {
			return storage.Object{}, err
		}
		if size != p.Size || hex.EncodeToString(hash.Sum(nil)) != p.ETag {
			return storage.Object{}, fmt.Errorf("分片 %d 大小或校验值不一致", p.Number)
		}
	}
	if err = ctx.Err(); err != nil {
		return storage.Object{}, err
	}
	if err = f.Sync(); err != nil {
		return storage.Object{}, err
	}
	if err = f.Close(); err != nil {
		return storage.Object{}, err
	}
	if err = s.root.Rename(temp, u.Key); err != nil {
		return storage.Object{}, err
	}
	obj, err := s.Stat(ctx, u.Key)
	if err != nil {
		return objectutil.Completed(u, ""), fmt.Errorf("文件已合并，读取元数据失败：%w", err)
	}
	// 合并成功后的清理错误同时返回文件信息，业务可再次取消会话清理残留。
	if err = s.cleanup(dir, u); err != nil {
		return obj, fmt.Errorf("文件已合并，清理分片失败：%w", err)
	}
	return obj, nil
}
func (s *Store) cleanup(dir string, u storage.MultipartUpload) error {
	for n := int64(1); n <= (u.Size-1)/u.PartSize+1; n++ {
		err := s.staging.Remove(dir + "/" + strconv.FormatInt(n, 10))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	f, err := s.staging.Open(dir)
	if err != nil {
		return err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return err
	}
	for _, name := range names {
		if strings.HasPrefix(name, "pending-") || strings.HasPrefix(name, "merged-") {
			if err = s.staging.Remove(dir + "/" + name); err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}
	}
	if err = s.staging.Remove(dir + "/session.json"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	err = s.staging.Remove(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
func (s *Store) AbortMultipart(ctx context.Context, u storage.MultipartUpload) error {
	s.multipartMu.Lock()
	defer s.multipartMu.Unlock()
	dir, err := s.load(ctx, u)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return s.cleanup(dir, u)
}
