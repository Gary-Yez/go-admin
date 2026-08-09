package sys_admin

import (
	"gitee.com/mxcker/go-admin/internal/state"
	"gitee.com/mxcker/go-admin/internal/system/sys_role"
)

func InitData() error {
	db := state.DB().Model(&SysAdmin{})
	count := int64(0)
	if err := db.Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		defaultRole := sys_role.SysRole{
			Default: true,
		}
		if err := state.DB().Model(sys_role.SysRole{}).Where(defaultRole).First(&defaultRole).Error; err != nil {
			panic(err)
		}
		data := []SysAdmin{
			{
				Username:     "admin",
				Nickname:     "默认管理员",
				Email:        "admin@qq.com",
				Phone:        "18888888888",
				Avatar:       "/img/user.png",
				PasswordHash: "$2a$10$PVIcAuZXvnP4sHLzGe/7se7F9Sakeu99ZwGqtlanUbFXgDHrxImQe",
				Default:      true,
				Role:         &defaultRole,
				RoleId:       defaultRole.Id,
			},
		}
		if err := db.Create(&data).Error; err != nil {
			return err
		}
	}
	return nil
}
