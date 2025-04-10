package sys_menu

import (
	"fmt"
	"gitee.com/mxcker/go-admin/server/core/global"
)

func InitData() error {
	fmt.Println("初始化系统菜单")
	data := []SysMenu{
		{
			Name:      "仪表盘",
			Icon:      "iconoir:arc-3d-center-point",
			Path:      "sys_home",
			Component: "../core/views/sys_home/index.vue",
			Sort:      0,
		},
		{
			Name: "系统设置",
			Icon: "iconoir:settings",
			Path: "sys_setting",
			Sort: 1,
			Children: []*SysMenu{
				{
					Name:      "菜单管理",
					Icon:      "iconoir:menu",
					Path:      "sys_menu",
					Component: "../core/views/sys_menu/index.vue",
					Sort:      0,
				}, {
					Name:      "角色管理",
					Icon:      "iconoir:user-badge-check",
					Path:      "sys_role",
					Component: "../core/views/sys_role/index.vue",
					Sort:      1,
				}, {
					Name:      "管理员管理",
					Icon:      "iconoir:user-crown",
					Path:      "sys_admin",
					Component: "../core/views/sys_admin/index.vue",
					Sort:      2,
				}, {
					Name:      "全局变量管理",
					Icon:      "iconoir:folder-settings",
					Path:      "sys_global_variable",
					Component: "../core/views/sys_global_variable/index.vue",
					Sort:      3,
				}, {
					Name:      "定时任务",
					Icon:      "iconoir:task-list",
					Path:      "sys_task",
					Component: "../core/views/sys_task/index.vue",
					Sort:      4,
				},
			},
		},
	}
	db := global.DB.Model(&SysMenu{})
	count := int64(0)
	if err := db.Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		if err := db.Create(&data).Error; err != nil {
			return err
		}
	}
	return nil
}
