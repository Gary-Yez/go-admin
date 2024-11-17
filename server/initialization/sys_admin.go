package initialization

import (
	"fmt"
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/models"
)

func InitSysAdmin() {
	fmt.Println("初始化默认管理员")
	db := global.DB.Model(&models.SysAdmin{})
	count := int64(0)
	if err := db.Count(&count).Error; err != nil {
		panic(err)
	}
	if count == 0 {
		defaultRole := &models.SysRole{}
		if err := global.DB.Model(&models.SysRole{}).First(defaultRole).Error; err != nil {
			panic(err)
		}
		data := []models.SysAdmin{
			{
				Username:     "admin",
				Nickname:     "默认管理员",
				Email:        "admin@qq.com",
				Phone:        "18888888888",
				Avatar:       "https://qmplusimg.henrongyi.top/gva_header.jpg",
				PasswordHash: "$2a$10$PVIcAuZXvnP4sHLzGe/7se7F9Sakeu99ZwGqtlanUbFXgDHrxImQe",
				Role:         defaultRole,
				RoleId:       defaultRole.Id,
			},
		}
		if err := db.Create(&data).Error; err != nil {
			panic(err)
		}
	}
}
