package sys_file

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/storage"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *serviceStruct) Upload(ctx context.Context, userId, storageId uint, name string, size int64, reader io.Reader, policy uploadPolicy) (*SysFile, error) {
	if err := validateExtension(name, policy.AllowedExtensions); err != nil {
		return nil, err
	}
	if size > policy.OrdinaryLimit {
		return nil, fmt.Errorf("超过 %d MiB 的文件请使用分片上传", policy.OrdinaryLimit>>20)
	}
	row, store, err := s.newFile(ctx, userId, storageId, name, size, false, 0, policy.SessionDays)
	if err != nil {
		return nil, err
	}
	defer closeStore(store)
	err = locked(ctx, row.Id, userId, func(tx *gorm.DB, file *SysFile) error {
		buffer := bufio.NewReader(reader)
		head, _ := buffer.Peek(512)
		file.ContentType = http.DetectContentType(head)
		if _, err := store.Put(ctx, file.ObjectKey, buffer, storage.PutOptions{Size: size, ContentType: file.ContentType}); err != nil {
			return err
		}
		if err := tx.Model(file).Update("content_type", file.ContentType).Error; err != nil {
			return err
		}
		if err := markReady(tx, file); err != nil {
			return err
		}
		row = file
		return nil
	})
	// 普通上传无法续传；独立清理上下文用于客户端断开后的善后。
	if err != nil {
		clean, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = s.Delete(clean, row.Id, userId, false)
	}
	return row, err
}

func (s *serviceStruct) Begin(ctx context.Context, userId uint, body *BeginBody) (*Session, error) {
	policy, err := readUploadPolicy()
	if err != nil {
		return nil, err
	}
	if err := validateExtension(body.Name, policy.AllowedExtensions); err != nil {
		return nil, err
	}
	if body.Size > policy.MaxSize {
		return nil, fmt.Errorf("文件过大，最多支持 %d 个分片", maxUploadParts)
	}
	if body.Size <= policy.OrdinaryLimit {
		return nil, fmt.Errorf("%d MiB 以内的文件请使用普通上传", policy.OrdinaryLimit>>20)
	}
	row, store, err := s.newFile(ctx, userId, body.StorageId, body.Name, body.Size, true, body.LastModified, policy.SessionDays)
	if err != nil {
		return nil, err
	}
	defer closeStore(store)
	var upload storage.MultipartUpload
	err = locked(ctx, row.Id, userId, func(tx *gorm.DB, file *SysFile) error {
		var err error
		upload, err = store.BeginMultipart(ctx, file.ObjectKey, storage.MultipartOptions{Size: file.Size, PartSize: policy.PartSize, ContentType: file.ContentType})
		if err != nil {
			return err
		}
		data, err := json.Marshal(upload)
		if err != nil {
			return err
		}
		return tx.Model(file).Update("upload_json", string(data)).Error
	})
	if err != nil {
		clean, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		// 已取得供应商会话时先保存再清理，失败仍可通过记录重试。
		if upload.ID != "" {
			data, _ := json.Marshal(upload)
			_ = state.DB().WithContext(clean).Model(row).Update("upload_json", string(data)).Error
			if abortErr := store.AbortMultipart(clean, upload); abortErr != nil {
				return nil, err
			}
		}
		_ = s.Delete(clean, row.Id, userId, false)
		return nil, err
	}
	row.IsOwner = true
	return &Session{File: row, Parts: []SysFilePart{}, PartSize: policy.PartSize}, nil
}

func (*serviceStruct) Part(ctx context.Context, id, userId uint, number int, reader io.Reader, length int64) (*SysFilePart, error) {
	var receipt *SysFilePart
	err := locked(ctx, id, userId, func(tx *gorm.DB, row *SysFile) error {
		if err := row.active(); err != nil {
			return err
		}
		if err := validateCurrentExtension(row.Name); err != nil {
			return err
		}
		upload, err := row.upload()
		if err != nil {
			return err
		}
		total := (upload.Size + upload.PartSize - 1) / upload.PartSize
		if number < 1 || int64(number) > total {
			return errors.New("分片编号无效")
		}
		expected := min(upload.PartSize, upload.Size-int64(number-1)*upload.PartSize)
		if length != expected {
			return errors.New("分片大小与会话不匹配")
		}
		store, err := row.openStore()
		if err != nil {
			return err
		}
		defer closeStore(store)
		buffer := bufio.NewReader(reader)
		if number == 1 {
			head, _ := buffer.Peek(512)
			// 原始供应商会话的 ContentType 不变，预览只使用服务端嗅探的类型。
			if err := tx.Model(row).Update("content_type", http.DetectContentType(head)).Error; err != nil {
				return err
			}
		}
		part, err := store.UploadPart(ctx, upload, number, buffer)
		if err != nil {
			return err
		}
		receipt = &SysFilePart{FileId: id, Number: part.Number, Size: part.Size, ETag: part.ETag}
		return tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "file_id"}, {Name: "number"}}, DoUpdates: clause.AssignmentColumns([]string{"size", "e_tag"})}).Create(receipt).Error
	})
	return receipt, err
}

func (*serviceStruct) Complete(ctx context.Context, id, userId uint) (*SysFile, error) {
	var result *SysFile
	err := locked(ctx, id, userId, func(tx *gorm.DB, row *SysFile) error {
		result = row
		if row.Status == ready {
			return nil
		}
		if err := row.active(); err != nil {
			return err
		}
		if err := validateCurrentExtension(row.Name); err != nil {
			return err
		}
		store, err := row.openStore()
		if err != nil {
			return err
		}
		defer closeStore(store)
		// 处理供应商合并成功但响应/数据库提交失败的重试。
		if object, err := store.Stat(ctx, row.ObjectKey); err == nil {
			if object.Size != row.Size {
				return errors.New("最终文件大小与会话不一致")
			}
			return markReady(tx, row)
		} else if !errors.Is(err, storage.ErrNotFound) {
			return err
		}
		upload, err := row.upload()
		if err != nil {
			return err
		}
		var saved []SysFilePart
		if err := tx.Where("file_id = ?", id).Order("number").Find(&saved).Error; err != nil {
			return err
		}
		if int64(len(saved)) != (upload.Size+upload.PartSize-1)/upload.PartSize {
			return errors.New("还有分片尚未上传，请继续上传")
		}
		parts := make([]storage.Part, 0, len(saved))
		for _, part := range saved {
			parts = append(parts, storage.Part{Number: part.Number, Size: part.Size, ETag: part.ETag})
		}
		if _, err := store.CompleteMultipart(ctx, upload, parts); err != nil {
			return err
		}
		return markReady(tx, row)
	})
	return result, err
}
