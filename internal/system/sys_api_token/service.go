package sys_api_token

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/Gary-Yez/go-admin/internal/cache"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_admin"
	"github.com/Gary-Yez/go-admin/internal/system/sys_auth"
	"github.com/Gary-Yez/go-admin/internal/system/sys_role"
	"github.com/Gary-Yez/go-admin/internal/utils"
	"github.com/Gary-Yez/go-admin/request"
	"gorm.io/gorm"
	"strings"
	"time"
)

type serviceStruct struct{}

func (*serviceStruct) List(ctx context.Context, req *request.ReqList) (list []TokenRow, total int64, err error) {
	// 派生表统一列名，避免连接后的 role_id、id 等字段歧义。
	rows := state.DB().WithContext(ctx).Model(&SysApiToken{}).Select(`sys_api_tokens.id, sys_api_tokens.admin_id, sys_api_tokens.role_id,
 sys_api_tokens.prefix, sys_api_tokens.remark, sys_api_tokens.created_at, sys_api_tokens.expires_at,
 COALESCE(a.username, '') AS username, COALESCE(a.nickname, '') AS nickname, COALESCE(r.name, '') AS role_name,
 CASE WHEN a.id IS NULL OR a.status <> 1 OR r.id IS NULL OR ar.admin_id IS NULL THEN 'invalid'
 WHEN sys_api_tokens.expires_at IS NOT NULL AND sys_api_tokens.expires_at <= ? THEN 'expired'
 WHEN sys_api_tokens.expires_at IS NULL THEN 'permanent' ELSE 'active' END AS status`, time.Now().UTC()).
		Joins("LEFT JOIN sys_admins a ON a.id = sys_api_tokens.admin_id").
		Joins("LEFT JOIN sys_roles r ON r.id = sys_api_tokens.role_id").
		Joins("LEFT JOIN sys_admin_role ar ON ar.admin_id = sys_api_tokens.admin_id AND ar.role_id = sys_api_tokens.role_id")
	db := req.WithFilter(state.DB().WithContext(ctx).Table("(?) AS tokens", rows), []string{"username", "remark", "role_id", "status", "prefix"})
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if req.Page == 0 {
		req.Page = 1
	}
	err = req.WithPagination(req.WithSort(db, []string{"id", "created_at"})).Order("id DESC").Scan(&list).Error
	return
}
func (*serviceStruct) RoleOptions(ctx context.Context) (list []sys_admin.RoleOption, err error) {
	err = state.DB().WithContext(ctx).Model(&sys_role.SysRole{}).Select("id", "name").Order("id").Scan(&list).Error
	return
}
func tokenHash(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *serviceStruct) SaveApiToken(adminId uint, body *ApiTokenBody) (string, error) {
	if body.ExpiresAt != nil && !body.ExpiresAt.After(time.Now()) {
		return "", errors.New("有效期必须晚于当前时间")
	}
	// 空有效期表示永久；有到期时间时统一转换为 UTC。
	var expiresAt *time.Time
	if body.ExpiresAt != nil {
		value := body.ExpiresAt.UTC()
		expiresAt = &value
	}
	if err := sys_auth.Service.VerifyAuthUser(&utils.AuthUser{UserId: adminId, RoleId: body.RoleId}); err != nil {
		return "", err
	}
	if body.Id != 0 {
		query := state.DB().Where("id = ? AND admin_id = ?", body.Id, adminId)
		err := mutateApiTokens(context.Background(), query, func(db *gorm.DB) error {
			result := db.Model(&SysApiToken{}).Where("id = ? AND admin_id = ?", body.Id, adminId).
				Updates(map[string]any{"role_id": body.RoleId, "expires_at": expiresAt, "remark": strings.TrimSpace(body.Remark)})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New("密钥不存在")
			}
			return nil
		})
		return "", err
	}

	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return "", err
	}
	token := "API_" + hex.EncodeToString(secret)
	row := SysApiToken{AdminId: adminId, RoleId: body.RoleId, TokenHash: tokenHash(token), Prefix: token[:12], ExpiresAt: expiresAt, Remark: strings.TrimSpace(body.Remark)}
	if err := state.DB().Create(&row).Error; err != nil {
		return "", err
	}
	return token, nil
}

// 所有删除入口共用缓存失效流程，避免全局管理绕过个人密钥缓存。
func (*serviceStruct) DeleteApiTokens(ctx context.Context, req *request.ReqIds, adminId uint) error {
	if err := req.Validate(); err != nil {
		return err
	}
	query := state.DB().Where("id IN ?", req.Ids)
	if adminId != 0 {
		query = query.Where("admin_id = ?", adminId)
	}
	return mutateApiTokens(ctx, query, func(db *gorm.DB) error {
		target := db.Where("id IN ?", req.Ids)
		if adminId != 0 {
			target = target.Where("admin_id = ?", adminId)
		}
		return target.Delete(&SysApiToken{}).Error
	})
}

type apiTokenIdentity struct {
	AdminId   uint       `json:"admin_id"`
	RoleId    uint       `json:"role_id"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func apiTokenCacheKey(hash string) string { return "auth:api-token:" + hash }

func lockApiToken(ctx context.Context, hash string) (cache.Lock, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	lock, err := state.Cache().WaitForLock(ctx, apiTokenCacheKey(hash), 30*time.Second, 10*time.Second)
	if err != nil {
		return nil, errors.New("密钥缓存繁忙或暂时不可用，请稍后重试")
	}
	return lock, nil
}

// 先持有回填共用的锁并清缓存，再提交数据库修改，防止旧数据并发回填。
func mutateApiTokens(ctx context.Context, query *gorm.DB, change func(*gorm.DB) error) error {
	var rows []SysApiToken
	if err := query.WithContext(ctx).Select("id", "token_hash").Order("token_hash").Find(&rows).Error; err != nil {
		return err
	}
	if len(rows) == 0 {
		return errors.New("密钥不存在")
	}
	locks := make([]cache.Lock, 0, len(rows))
	defer func() {
		for i := len(locks) - 1; i >= 0; i-- {
			_ = locks[i].Unlock()
		}
	}()
	for _, row := range rows {
		lock, err := lockApiToken(ctx, row.TokenHash)
		if err != nil {
			return err
		}
		locks = append(locks, lock)
	}
	for _, row := range rows {
		if err := state.Cache().Del(apiTokenCacheKey(row.TokenHash)); err != nil {
			return errors.New("清除密钥缓存失败，请重试")
		}
	}
	return state.DB().WithContext(ctx).Transaction(change)
}

func readApiToken(ctx context.Context, hash string) (*apiTokenIdentity, error) {
	store := state.Cache()
	key := apiTokenCacheKey(hash)
	var identity apiTokenIdentity
	err := store.GetJSON(key, &identity)
	if err == nil {
		return &identity, nil
	}
	if !errors.Is(err, cache.ErrCacheNotFound) {
		return nil, errors.New("密钥缓存暂时不可用")
	}
	lock, err := lockApiToken(ctx, hash)
	if err != nil {
		return nil, err
	}
	defer lock.Unlock()
	if err = store.GetJSON(key, &identity); err == nil {
		return &identity, nil
	}
	if !errors.Is(err, cache.ErrCacheNotFound) {
		return nil, errors.New("密钥缓存暂时不可用")
	}
	var row SysApiToken
	err = state.DB().WithContext(ctx).Select("admin_id", "role_id", "expires_at").Where("token_hash = ? AND (expires_at IS NULL OR expires_at > ?)", hash, time.Now().UTC()).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("API 密钥无效、已过期或已删除")
		}
		return nil, errors.New("API 密钥验证失败")
	}
	identity = apiTokenIdentity{AdminId: row.AdminId, RoleId: row.RoleId, ExpiresAt: row.ExpiresAt}
	ttl := cache.DefaultTTL
	if row.ExpiresAt != nil {
		if remaining := time.Until(*row.ExpiresAt); remaining < ttl {
			ttl = remaining
		}
	}
	if ttl <= 0 {
		return nil, errors.New("API 密钥已过期")
	}
	if err := store.SetJSON(key, identity, ttl); err != nil {
		return nil, errors.New("密钥缓存暂时不可用")
	}
	return &identity, nil
}

func (s *serviceStruct) VerifyApiToken(ctx context.Context, token string) (*utils.AuthUser, error) {
	row, err := readApiToken(ctx, tokenHash(token))
	if err != nil {
		return nil, err
	}
	if row.ExpiresAt != nil && !row.ExpiresAt.After(time.Now()) {
		return nil, errors.New("API 密钥已过期")
	}
	return &utils.AuthUser{UserId: row.AdminId, RoleId: row.RoleId}, nil
}
