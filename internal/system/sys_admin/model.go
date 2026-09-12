package sys_admin

import (
	"github.com/Gary-Yez/go-admin/internal/system/sys_role"
	"github.com/Gary-Yez/go-admin/request"
	"time"
)

type SysAdmin struct {
	Id           uint                `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt    time.Time           `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt    time.Time           `json:"updated_at" gorm:"comment:更新时间"`
	Username     string              `json:"username" gorm:"unique;comment:用户名"`
	Avatar       string              `json:"avatar" gorm:"comment:头像"`
	Nickname     string              `json:"nickname" gorm:"comment:昵称"`
	Email        string              `json:"email" gorm:"unique;comment:邮箱"`
	Phone        string              `json:"phone" gorm:"unique;comment:手机号"`
	Password     string              `json:"password" gorm:"-"`
	PasswordHash string              `json:"-"`
	Status       uint                `json:"status" gorm:"default:1;comment:状态"`
	RoleId       uint                `json:"role_id"`
	Role         *sys_role.SysRole   `json:"role" gorm:"foreignKey:role_id;"`
	Roles        []*sys_role.SysRole `json:"roles" gorm:"many2many:sys_admin_role;joinForeignKey:AdminId;joinReferences:RoleId;"`
	RoleIds      []uint              `json:"role_ids" gorm:"-"`
	LoginVersion uint64              `json:"-" gorm:"not null;default:1;comment:登录版本"`
}

// RoleId 表示登录的默认角色，当前角色由登录令牌携带。
type SysAdminRole struct {
	AdminId uint `gorm:"primaryKey"`
	RoleId  uint `gorm:"primaryKey"`
}

func (SysAdminRole) TableName() string { return "sys_admin_role" }

type RoleOption struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

// RoleId 按所有已绑定角色筛选，而非账号的默认角色。
type AdminListQuery struct {
	request.ReqList
	RoleId uint `json:"role_id" form:"role_id"`
}
