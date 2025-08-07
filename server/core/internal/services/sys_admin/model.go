package sys_admin

import (
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_role"
	"gitee.com/mxcker/go-admin/server/utils"
)

type SysAdmin struct {
	utils.DbBaseModel
	Username     string            `json:"username" gorm:"unique;comment:用户名"`
	Avatar       string            `json:"avatar" gorm:"comment:头像"`
	Nickname     string            `json:"nickname" gorm:"comment:昵称"`
	Email        string            `json:"email" gorm:"unique;comment:邮箱"`
	Phone        string            `json:"phone" gorm:"unique;comment:手机号"`
	Password     string            `json:"password" gorm:"-"`
	PasswordHash string            `json:"-"`
	Status       uint              `json:"status" gorm:"default:1;comment:状态"`
	Default      bool              `json:"is_default"`
	ApiToken     string            `json:"api_token" gorm:"unique;default:null;comment:Api密钥"`
	RoleId       uint              `json:"role_id"`
	Role         *sys_role.SysRole `json:"role" gorm:"foreignKey:role_id;"`
}
