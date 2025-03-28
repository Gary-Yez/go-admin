package sys_menu

import (
	"fmt"
	"gitee.com/mxcker/go-admin/server/core/global"
)

func Init() {
	fmt.Println("初始化系统菜单")
	data := []SysMenu{
		{
			Name:      "仪表盘",
			Icon:      "Odometer",
			Path:      "sys_home",
			Component: "../views/sys_home/index.vue",
			Sort:      0,
		},
		{
			Name: "系统设置",
			Icon: "Setting",
			Path: "sys_setting",
			Sort: 1,
			Children: []*SysMenu{
				{
					Name:      "菜单管理",
					Icon:      "Menu",
					Path:      "sys_menu",
					Component: "../views/sys_menu/index.vue",
					Sort:      0,
				}, {
					Name:      "角色管理",
					Icon:      "UserFilled",
					Path:      "sys_role",
					Component: "../views/sys_role/index.vue",
					Sort:      1,
				}, {
					Name:      "管理员管理",
					Icon:      "User",
					Path:      "sys_admin",
					Component: "../views/sys_admin/index.vue",
					Sort:      2,
				},
			},
		},
	}
	db := global.DB.Model(&SysMenu{})
	count := int64(0)
	if err := db.Count(&count).Error; err != nil {
		panic(err)
	}
	if count == 0 {
		if err := db.Create(&data).Error; err != nil {
			panic(err)
		}
	}
}
