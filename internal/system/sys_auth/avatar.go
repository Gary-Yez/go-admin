package sys_auth

import (
	"bytes"
	"context"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/Gary-Yez/go-admin/internal/system/sys_file"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const maxAvatarSize = 5 << 20

// validateAvatar 只校验文件并生成预览地址，资料更新由 ChangeInfo 的事务统一完成。
func validateAvatar(ctx context.Context, tx *gorm.DB, userId, fileId uint) (string, error) {
	var file sys_file.SysFile
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND user_id = ? AND status = ?", fileId, userId, "ready").First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("请选择本人上传且已完成的文件")
		}
		return "", err
	}
	if file.Size <= 0 || file.Size > maxAvatarSize {
		return "", errors.New("头像大小不能超过 5 MiB")
	}
	if err := file.ReadContent(ctx, func(reader io.Reader) error {
		data, err := io.ReadAll(io.LimitReader(reader, maxAvatarSize+1))
		if err != nil {
			return err
		}
		if len(data) > maxAvatarSize {
			return errors.New("头像大小不能超过 5 MiB")
		}
		config, format, err := image.DecodeConfig(bytes.NewReader(data))
		if err != nil || (format != "jpeg" && format != "png") {
			return errors.New("头像必须为有效的 JPG 或 PNG 图片")
		}
		if config.Width <= 0 || config.Height <= 0 || config.Width > 4096 || config.Height > 4096 {
			return errors.New("头像宽高不能超过 4096 像素")
		}
		if _, _, err := image.Decode(bytes.NewReader(data)); err != nil {
			return errors.New("头像图片已损坏，请重新上传")
		}
		return nil
	}); err != nil {
		return "", err
	}
	return file.PreviewLink(ctx)
}
