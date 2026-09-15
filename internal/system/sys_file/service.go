package sys_file

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_storage"
	"github.com/Gary-Yez/go-admin/request"
	"github.com/Gary-Yez/go-admin/storage"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type serviceStruct struct{}

func (*serviceStruct) List(ctx context.Context, req *request.ReqList, userId uint) (list []SysFile, total int64, err error) {
	db := req.WithFilter(state.DB().WithContext(ctx).Model(&SysFile{}), []string{"name", "engine", "status", "username", "created_at", "storage_id"})
	if err = db.Count(&total).Error; err != nil {
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	err = req.WithPagination(req.WithSort(db, []string{"id", "size", "created_at", "updated_at"})).Order("id DESC").Omit("UploadJSON", "ObjectKey").Find(&list).Error
	if err != nil || len(list) == 0 {
		return
	}
	ids := make([]uint, 0, len(list))
	for _, row := range list {
		ids = append(ids, row.StorageId)
	}
	var accounts []sys_storage.SysStorage
	if err = state.DB().WithContext(ctx).Select("id", "name").Where("id IN ?", ids).Find(&accounts).Error; err != nil {
		return
	}
	names := map[uint]string{}
	for _, account := range accounts {
		names[account.Id] = account.Name
	}
	for i := range list {
		list[i].StorageName = names[list[i].StorageId]
		list[i].IsOwner = list[i].UserId == userId
	}
	return
}

func randomKey() (string, error) {
	var value [32]byte
	if _, err := rand.Read(value[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(value[:]), nil
}

func (s *serviceStruct) newFile(ctx context.Context, userId, storageId uint, name string, size int64, multipart bool, lastModified, sessionDays int64) (*SysFile, storage.Store, error) {
	name = filepath.Base(strings.ReplaceAll(name, "\\", "/"))
	if strings.TrimSpace(name) == "" || name == "." || len([]rune(name)) > 255 || strings.ContainsAny(name, "\r\n\x00") || size < 0 {
		return nil, nil, errors.New("文件名称或大小不符合要求")
	}
	key, err := randomKey()
	if err != nil {
		return nil, nil, err
	}
	expires := time.Now().Add(time.Duration(sessionDays) * 24 * time.Hour)
	row := &SysFile{Name: name, Size: size, ContentType: "application/octet-stream", Status: uploading, UserId: userId, Multipart: multipart,
		LastModified: lastModified, ObjectKey: time.Now().Format("2006/01/02") + "/" + key, ExpiresAt: &expires}
	var store storage.Store
	err = sys_storage.WithUploadAccount(ctx, storageId, func(tx *gorm.DB, account *sys_storage.SysStorage, config storage.Config) error {
		var err error
		store, err = storage.New(config)
		if err != nil {
			return err
		}
		row.StorageId = account.Id
		row.StorageName = account.Name
		row.Engine = string(account.Engine)
		if err := tx.Table("sys_admins").Select("username").Where("id = ?", userId).Scan(&row.Username).Error; err != nil {
			return err
		}
		return tx.Create(row).Error
	})
	if err != nil {
		if store != nil {
			closeStore(store)
		}
		return nil, nil, err
	}
	return row, store, nil
}

func closeStore(store storage.Store) {
	if closer, ok := store.(io.Closer); ok {
		_ = closer.Close()
	}
}
func (f *SysFile) openStore() (storage.Store, error) {
	config, err := sys_storage.ReadConfig(f.StorageId, false)
	if err != nil {
		return nil, err
	}
	return storage.New(config)
}
func (f *SysFile) upload() (storage.MultipartUpload, error) {
	var upload storage.MultipartUpload
	if f.UploadJSON == "" {
		return upload, errors.New("上传会话未建立，请删除记录后重新上传")
	}
	err := json.Unmarshal([]byte(f.UploadJSON), &upload)
	return upload, err
}

// 会话行锁协调分片、合并和删除；所有实例使用同一数据库。
// 存储操作不受数据库回滚保护，保留记录供重试，合并重试先检查最终对象。
func locked(ctx context.Context, id, owner uint, fn func(*gorm.DB, *SysFile) error) error {
	return state.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row SysFile
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&row, id).Error
		if err != nil {
			return err
		}
		if owner != 0 && row.UserId != owner {
			return errors.New("只能操作自己创建的上传会话")
		}
		return fn(tx, &row)
	})
}
func (f *SysFile) active() error {
	if f.Status != uploading {
		return errors.New("文件已上传完成")
	}
	if f.ExpiresAt == nil || !f.ExpiresAt.After(time.Now()) {
		return errors.New("上传会话已过期，请删除后重新上传")
	}
	return nil
}
func markReady(tx *gorm.DB, row *SysFile) error {
	if err := tx.Model(row).Updates(map[string]any{"status": ready, "expires_at": nil, "upload_json": ""}).Error; err != nil {
		return err
	}
	row.Status = ready
	row.ExpiresAt = nil
	row.UploadJSON = ""
	if row.Multipart {
		return tx.Where("file_id = ?", row.Id).Delete(&SysFilePart{}).Error
	}
	return nil
}

func (*serviceStruct) Session(ctx context.Context, id, owner uint) (*Session, error) {
	result := &Session{Parts: []SysFilePart{}}
	err := locked(ctx, id, owner, func(tx *gorm.DB, row *SysFile) error {
		if row.Status != ready {
			if err := row.active(); err != nil {
				return err
			}
		}
		row.IsOwner = true
		result.File = row
		if row.Multipart && row.Status != ready {
			upload, err := row.upload()
			if err != nil {
				return err
			}
			result.PartSize = upload.PartSize
		}
		return tx.Where("file_id = ?", id).Order("number").Find(&result.Parts).Error
	})
	return result, err
}

// Delete 默认先删除对象/分片，再删除记录；recordsOnly 仅供文件管理清理记录，不访问存储。
func (*serviceStruct) Delete(ctx context.Context, id, owner uint, recordsOnly bool) error {
	err := locked(ctx, id, owner, func(tx *gorm.DB, row *SysFile) error {
		if owner != 0 && row.Status == ready {
			return errors.New("已完成的文件请通过删除接口管理")
		}
		if recordsOnly && owner == 0 {
			return deleteRecords(tx, row)
		}
		return deleteFile(ctx, tx, row)
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func deleteFile(ctx context.Context, tx *gorm.DB, row *SysFile) error {
	store, err := row.openStore()
	if err != nil {
		return err
	}
	defer closeStore(store)
	if row.UploadJSON != "" {
		upload, err := row.upload()
		if err != nil {
			return err
		}
		if err := store.AbortMultipart(ctx, upload); err != nil {
			return err
		}
	}
	if err := store.Delete(ctx, row.ObjectKey); err != nil {
		return err
	}
	return deleteRecords(tx, row)
}

func deleteRecords(tx *gorm.DB, row *SysFile) error {
	if err := tx.Where("file_id = ?", row.Id).Delete(&SysFilePart{}).Error; err != nil {
		return err
	}
	return tx.Delete(row).Error
}

func (*serviceStruct) Cleanup(ctx context.Context) (int, error) {
	var ids []uint
	if err := state.DB().WithContext(ctx).Model(&SysFile{}).Where("status = ? AND expires_at < ?", uploading, time.Now()).Pluck("id", &ids).Error; err != nil {
		return 0, err
	}
	count := 0
	var failures []error
	for _, id := range ids {
		if err := ctx.Err(); err != nil {
			failures = append(failures, err)
			break
		}
		deleted := false
		err := locked(ctx, id, 0, func(tx *gorm.DB, row *SysFile) error {
			// 等待行锁期间上传可能已经完成，取得锁后再次核对，不能误删成品。
			if row.Status != uploading || row.ExpiresAt == nil || row.ExpiresAt.After(time.Now()) {
				return nil
			}
			if err := deleteFile(ctx, tx, row); err != nil {
				return err
			}
			deleted = true
			return nil
		})
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				failures = append(failures, fmt.Errorf("会话 %d 清理失败：%w", id, err))
			}
			continue
		}
		if deleted {
			count++
		}
	}
	if err := errors.Join(failures...); err != nil {
		return count, fmt.Errorf("已清理 %d 条，部分会话未完成清理：%w", count, err)
	}
	return count, nil
}
