package sys_auth

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/Gary-Yez/go-admin/dberror"
	"github.com/Gary-Yez/go-admin/internal/cache"
	"github.com/Gary-Yez/go-admin/internal/state"
	utils2 "github.com/Gary-Yez/go-admin/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"

	"github.com/Gary-Yez/go-admin/internal/system/sys_admin"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"
	"gorm.io/gorm"
	"slices"
)

type serviceStruct struct {
}

var errLoginLocked = errors.New("当前来源登录此账号失败次数过多，请在锁定结束后重试")

// 失败统计以账号和连接来源为单位，通过共享缓存锁串行更新。
type loginFailures struct {
	Count       int       `json:"count"`
	WindowEnd   time.Time `json:"window_end"`
	LockedUntil time.Time `json:"locked_until"`
}

func (s *serviceStruct) Login(ctx context.Context, username, password, ip string) (*sys_admin.SysAdmin, error) {
	maximum, err := utils2.ReadLoginInt("login.max_failures")
	if err != nil {
		return nil, err
	}
	window, err := utils2.ReadLoginInt("login.failure_window_minutes")
	if err != nil {
		return nil, err
	}
	duration, err := utils2.ReadLoginInt("login.lock_minutes")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	admin := new(sys_admin.SysAdmin)
	lookupErr := state.DB().WithContext(ctx).Where("username = ?", username).First(admin).Error
	if lookupErr != nil && !errors.Is(lookupErr, gorm.ErrRecordNotFound) {
		return nil, errors.New("登录暂时不可用，请稍后重试")
	}
	// 使用数据库中的账号名，避免大小写等数据库排序规则差异绕过计数。
	identity := strings.ToLower(strings.TrimSpace(username))
	if lookupErr == nil {
		identity = admin.Username
	}
	key := fmt.Sprintf("auth:login-failures:%x", sha256.Sum256([]byte(identity+"\x00"+ip)))
	store := state.Cache()
	lock, err := store.Lock(ctx, key, 30*time.Second, 10*time.Second)
	if err != nil {
		return nil, errors.New("登录请求繁忙，请稍后重试")
	}
	defer lock.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, errors.New("登录请求超时，请重试")
	}
	var failures loginFailures
	if err := store.GetJSON(key, &failures); err != nil && !errors.Is(err, cache.ErrCacheNotFound) {
		return nil, errors.New("登录限制暂时不可用，请稍后重试")
	}
	now := time.Now()
	if now.Before(failures.LockedUntil) {
		return nil, errLoginLocked
	}
	if !failures.LockedUntil.IsZero() || !now.Before(failures.WindowEnd) {
		failures = loginFailures{WindowEnd: now.Add(time.Duration(window) * time.Minute)}
	}
	hash := admin.PasswordHash
	if lookupErr != nil {
		// 不存在的账号也进行密码计算并累计失败，返回相同提示。
		hash = "$2a$10$PVIcAuZXvnP4sHLzGe/7se7F9Sakeu99ZwGqtlanUbFXgDHrxImQe"
	}
	passwordErr := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if lookupErr != nil || passwordErr != nil {
		if !time.Now().Before(failures.WindowEnd) {
			failures = loginFailures{WindowEnd: time.Now().Add(time.Duration(window) * time.Minute)}
		}
		failures.Count++
		expires := failures.WindowEnd
		if failures.Count >= maximum {
			failures.LockedUntil = time.Now().Add(time.Duration(duration) * time.Minute)
			expires = failures.LockedUntil
		}
		if err := store.SetJSON(key, failures, time.Until(expires)); err != nil {
			return nil, errors.New("登录限制暂时不可用，请稍后重试")
		}
		if !failures.LockedUntil.IsZero() {
			return nil, errLoginLocked
		}
		return nil, errors.New("用户名或密码不正确")
	}
	if err := s.VerifyAuthUser(&utils2.AuthUser{UserId: admin.Id, RoleId: admin.RoleId, LoginVersion: &admin.LoginVersion}); err != nil {
		return nil, err
	}
	if err := store.Del(key); err != nil {
		return nil, errors.New("清除登录失败记录失败，请稍后重试")
	}
	return admin, nil
}

func (s *serviceStruct) GetUser(id uint, roleId uint) (user *sys_admin.SysAdmin, err error) {
	user = &sys_admin.SysAdmin{}
	err = state.DB().Preload("Roles").Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}
	if user.Status != 1 {
		return nil, errors.New("管理员被禁用")
	}
	user.RoleIds = make([]uint, 0, len(user.Roles))
	for _, role := range user.Roles {
		user.RoleIds = append(user.RoleIds, role.Id)
		if role.Id == roleId {
			user.Role = role
		}
	}
	if user.Role == nil {
		return nil, errors.New("当前角色已被移除，请重新登录")
	}
	user.RoleId = roleId
	var menus []*sys_menu.SysMenu
	if user.Role.IsSuperAdmin {
		if err = state.DB().Model(&sys_menu.SysMenu{}).Find(&menus).Error; err != nil {
			return nil, err
		}
		user.Role.Menus = menus
	} else {
		if err := state.DB().Model(user.Role).Association("Menus").Find(&menus); err != nil {
			return nil, err
		}
		user.Role.Menus, err = sys_menu.Service.WithAncestors(menus)
		if err != nil {
			return nil, err
		}
	}
	return user, err
}

// 对 JWT 和 API 密钥统一核对账号状态与角色归属，不能只信任令牌中的角色ID。
func (s *serviceStruct) VerifyAuthUser(user *utils2.AuthUser) error {
	if user == nil || user.UserId == 0 || user.RoleId == 0 {
		return errors.New("登录身份无效")
	}
	identity, err := sys_admin.Service.ReadIdentity(user.UserId)
	if err != nil {
		return err
	}
	if user.LoginVersion != nil && *user.LoginVersion != identity.LoginVersion {
		return errors.New("登录已失效，请重新登录")
	}
	if identity.Status != 1 {
		return errors.New("管理员被禁用")
	}
	if !slices.Contains(identity.RoleIds, user.RoleId) {
		return errors.New("当前角色已被移除，请重新登录")
	}
	return nil
}

// ChangeInfo 在一个事务中保存提交的字段，头像校验失败时不写入其他资料。
func (*serviceStruct) ChangeInfo(ctx context.Context, userId uint, body *ChangeInfoBody) (map[string]any, error) {
	updates := make(map[string]any)
	for _, field := range []struct {
		key, label string
		value      *string
	}{
		{"nickname", "昵称", body.Nickname},
		{"phone", "手机号", body.Phone},
		{"email", "邮箱", body.Email},
	} {
		if field.value == nil {
			continue
		}
		value := strings.TrimSpace(*field.value)
		if value == "" {
			return nil, fmt.Errorf("%s不能为空", field.label)
		}
		updates[field.key] = value
	}
	if body.AvatarFileId != nil {
		if *body.AvatarFileId == 0 {
			return nil, errors.New("请选择已上传的头像")
		}
		updates["avatar_file_id"] = *body.AvatarFileId
	}
	if len(updates) == 0 {
		return nil, errors.New("请提交需要修改的资料")
	}
	var avatar string
	err := state.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if body.AvatarFileId != nil {
			var err error
			avatar, err = validateAvatar(ctx, tx, userId, *body.AvatarFileId)
			if err != nil {
				return err
			}
		}
		return dberror.Unique(tx.Model(&sys_admin.SysAdmin{}).Where("id = ?", userId).Updates(updates).Error, &sys_admin.SysAdmin{})
	})
	if err != nil {
		return nil, err
	}
	if body.AvatarFileId != nil {
		updates["avatar"] = avatar
	}
	return updates, nil
}
