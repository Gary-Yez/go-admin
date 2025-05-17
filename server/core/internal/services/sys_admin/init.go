package sys_admin

import (
	"fmt"
	"gitee.com/mxcker/go-admin/server/core/internal/services/sys_role"
	"gitee.com/mxcker/go-admin/server/global"
)

func InitData() error {
	fmt.Println("初始化默认管理员")
	db := global.DB.Model(&SysAdmin{})
	count := int64(0)
	if err := db.Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		defaultRole := sys_role.SysRole{
			Default: true,
		}
		if err := global.DB.Model(sys_role.SysRole{}).Where(defaultRole).First(&defaultRole).Error; err != nil {
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
