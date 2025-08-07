package sys_menu

import (
	"gitee.com/mxcker/go-admin/server/global"
)

func InitData() error {
	data := []*SysMenu{
		{
			Name:      "仪表盘",
			Key:       "sys_home",
			Icon:      "iconoir:arc-3d-center-point",
			Path:      "sys_home",
			Component: "../core/views/sys_home/index.vue",
			Sort:      0,
		},
		{
			Name: "系统设置",
			Key:  "sys_setting",
			Icon: "iconoir:settings",
			Path: "sys_setting",
			Sort: 1,
			Children: []*SysMenu{
				{
					Name:      "菜单管理",
					Key:       "sys_menu",
					Icon:      "iconoir:menu",
					Path:      "sys_menu",
					Component: "../core/views/sys_menu/index.vue",
					Sort:      0,
				}, {
					Name:      "API管理",
					Key:       "sys_apis",
					Icon:      "iconoir:laptop-charging",
					Path:      "sys_apis",
					Component: "../core/views/sys_apis/index.vue",
					Sort:      1,
				}, {
					Name:      "角色管理",
					Key:       "sys_role",
					Icon:      "iconoir:user-badge-check",
					Path:      "sys_role",
					Component: "../core/views/sys_role/index.vue",
					Sort:      2,
				}, {
					Name:      "管理员管理",
					Key:       "sys_admin",
					Icon:      "iconoir:user-crown",
					Path:      "sys_admin",
					Component: "../core/views/sys_admin/index.vue",
					Sort:      3,
				}, {
					Name:      "计划任务",
					Key:       "sys_cron_job",
					Icon:      "iconoir:task-list",
					Path:      "sys_cron_job",
					Component: "../core/views/sys_cron_job/index.vue",
					Sort:      4,
				},
			},
		}, {
			Name:      "个人信息",
			Key:       "sys_userinfo",
			Icon:      "iconoir:arc-3d-center-point",
			Path:      "sys_userinfo",
			Component: "../core/views/sys_userinfo/index.vue",
			Sort:      2,
			Hidden:    true,
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
