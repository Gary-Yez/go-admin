package initialization

import (
	"fmt"
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/models"
)

func InitSysMenu() {
	fmt.Println("初始化系统菜单")
	data := []models.SysMenu{
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
			Children: []*models.SysMenu{
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
		{
			Name:      "开发工具",
			Icon:      "EditPen",
			Path:      "sys_dev_tools",
			Component: "../views/sys_home/index.vue",
			Sort:      2,
			Children: []*models.SysMenu{
				{
					Name:      "代码生成",
					Icon:      "Cpu",
					Path:      "sys_autocode",
					Component: "../views/sys_autocode/index.vue",
					Sort:      0,
				}, {
					Name:      "生成历史",
					Icon:      "Document",
					Path:      "sys_autocode_history",
					Component: "../views/sys_autocode/history.vue",
					Sort:      1,
				},
			},
		},
	}
	db := global.DB.Model(&models.SysMenu{})
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
