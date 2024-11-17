package initialization

import (
	"fmt"
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/models"
)

func InitSysRole() {
	fmt.Println("初始化默认角色")
	db := global.DB.Model(&models.SysRole{})
	count := int64(0)
	if err := db.Count(&count).Error; err != nil {
		panic(err)
	}
	if count == 0 {
		var menuList []*models.SysMenu
		if err := global.DB.Model(&models.SysMenu{}).Find(&menuList).Error; err != nil {
			panic(err)
		}
		data := []models.SysRole{
			{
				Name:  "默认管理员",
				Menus: menuList,
			},
		}
		if err := db.Create(&data).Error; err != nil {
			panic(err)
		}
	}
}
