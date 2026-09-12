package sys_api_token

import "time"

type SysApiToken struct {
	Id        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	AdminId   uint       `json:"-" gorm:"index;not null"`
	RoleId    uint       `json:"role_id" gorm:"not null"`
	TokenHash string     `json:"-" gorm:"size:64;uniqueIndex;not null"`
	Prefix    string     `json:"prefix" gorm:"size:16;not null"`
	ExpiresAt *time.Time `json:"expires_at"`
	Remark    string     `json:"remark" gorm:"size:200"`
}

type ApiTokenBody struct {
	Id        uint       `json:"id"`
	RoleId    uint       `json:"role_id" binding:"required"`
	ExpiresAt *time.Time `json:"expires_at"`
	Remark    string     `json:"remark" binding:"max=200"`
}

// 显式投影，管理接口不读取或返回密钥哈希。
type TokenRow struct {
	Id        uint       `json:"id"`
	AdminId   uint       `json:"admin_id"`
	Username  string     `json:"username"`
	Nickname  string     `json:"nickname"`
	RoleId    uint       `json:"role_id"`
	RoleName  string     `json:"role_name"`
	Prefix    string     `json:"prefix"`
	Remark    string     `json:"remark"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at"`
	Status    string     `json:"status"`
}
