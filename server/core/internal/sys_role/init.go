package sys_role

import (
	"fmt"
	"gitee.com/mxcker/go-admin/server/core/global"
)

func Init() {
	fmt.Println("初始化默认角色")
	db := global.DB.Model(&SysRole{})
	count := int64(0)
	if err := db.Count(&count).Error; err != nil {
		panic(err)
	}
	if count == 0 {
		data := []SysRole{
			{
				Name:    "超级管理员",
				Default: true,
			},
		}
		if err := db.Create(&data).Error; err != nil {
			panic(err)
		}
	}
}
