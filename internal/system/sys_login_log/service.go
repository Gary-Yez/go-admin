package sys_login_log

import (
	"context"
	"errors"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/request"
	"strings"
	"time"
	"unicode/utf8"
)

type serviceStruct struct{}

// Record 独立提交，客户端断开也尝试记录；不接收密码、令牌或请求正文。
func (*serviceStruct) Record(row *SysLoginLog) error {
	row.Username = boundedText(row.Username, 191)
	row.IP = boundedText(row.IP, 45)
	row.UserAgent = boundedText(row.UserAgent, 512)
	row.Message = boundedText(row.Message, 512)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return state.DB().WithContext(ctx).Create(row).Error
}

func boundedText(value string, limit int) string {
	value = strings.ToValidUTF8(value, "")
	if len(value) <= limit {
		return value
	}
	value = value[:limit]
	for !utf8.ValidString(value) {
		value = value[:len(value)-1]
	}
	return value
}

func (*serviceStruct) List(ctx context.Context, req *request.ReqList) (list []*SysLoginLog, total int64, err error) {
	db := req.WithFilter(state.DB().WithContext(ctx).Model(&SysLoginLog{}), []string{"username", "ip", "status", "created_at"})
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if req.Page == 0 {
		req.Page = 1
	}
	err = req.WithPagination(req.WithSort(db, []string{"id", "created_at", "duration"})).Order("id DESC").Find(&list).Error
	return
}

func (*serviceStruct) Delete(ctx context.Context, req *request.ReqIds) error {
	return req.WithQuery(state.DB().WithContext(ctx)).Delete(&SysLoginLog{}).Error
}

func (*serviceStruct) Cleanup(ctx context.Context, days int) (int64, error) {
	if days < 1 || days > 36500 {
		return 0, errors.New("保留天数必须在 1–36500 之间")
	}
	result := state.DB().WithContext(ctx).Where("created_at < ?", time.Now().AddDate(0, 0, -days)).Delete(&SysLoginLog{})
	return result.RowsAffected, result.Error
}
