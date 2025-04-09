package sys_admin

import (
	"gitee.com/mxcker/go-admin/server/core/internal/sys_role"
	"gitee.com/mxcker/go-admin/server/core/types"
)

type SysAdmin struct {
	types.DbBaseModel
	Username     string            `json:"username" gorm:"unique;comment:用户名"`
	Avatar       string            `json:"avatar" gorm:"comment:头像"`
	Nickname     string            `json:"nickname" gorm:"comment:昵称"`
	Email        string            `json:"email" gorm:"unique;comment:邮箱"`
	Phone        string            `json:"phone" gorm:"unique;comment:手机号"`
	Password     string            `json:"password" gorm:"-"`
	PasswordHash string            `json:"-"`
	Status       uint              `json:"status" gorm:"default:1;comment:状态"`
	Default      bool              `json:"is_default"`
	RoleId       uint              `json:"role_id"`
	Role         *sys_role.SysRole `json:"role" gorm:"foreignKey:role_id;"`
}
