package sys_apis

import (
	"gitee.com/mxcker/go-admin/server/global"
	"net/http"
)

func InitData() error {
	count := int64(0)
	if err := global.DB.Model(SysApi{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		var defaultApis = []SysApi{
			// 系统模块-鉴权
			{Method: http.MethodGet, Path: "/sys_auth/me", Group: "系统模块-鉴权", Description: "获取当前用户(必选)"},
			{Method: http.MethodPost, Path: "/sys_auth/change_info", Group: "系统模块-鉴权", Description: "修改自身信息"},
			{Method: http.MethodPost, Path: "/sys_auth/change_password", Group: "系统模块-鉴权", Description: "修改自身密码"},
			// 系统模块-管理员
			{Method: http.MethodPost, Path: "/sys_admin/create", Group: "系统模块-管理员", Description: "创建管理员"},
			{Method: http.MethodPost, Path: "/sys_admin/delete", Group: "系统模块-管理员", Description: "删除管理员"},
			{Method: http.MethodPost, Path: "/sys_admin/edit", Group: "系统模块-管理员", Description: "修改管理员"},
			{Method: http.MethodGet, Path: "/sys_admin/get", Group: "系统模块-管理员", Description: "查询管理员"},
			{Method: http.MethodGet, Path: "/sys_admin/list", Group: "系统模块-管理员", Description: "管理员列表"},
			// 系统模块-API
			{Method: http.MethodPost, Path: "/sys_apis/create", Group: "系统模块-API", Description: "创建API"},
			{Method: http.MethodPost, Path: "/sys_apis/delete", Group: "系统模块-API", Description: "删除API"},
			{Method: http.MethodPost, Path: "/sys_apis/edit", Group: "系统模块-API", Description: "修改API"},
			{Method: http.MethodGet, Path: "/sys_apis/get", Group: "系统模块-API", Description: "获取API"},
			{Method: http.MethodGet, Path: "/sys_apis/get_groups", Group: "系统模块-API", Description: "API分组列表"},
			{Method: http.MethodPost, Path: "/sys_apis/list", Group: "系统模块-API", Description: "API列表"},
			{Method: http.MethodGet, Path: "/sys_apis/sync_api", Group: "系统模块-API", Description: "同步API"},
			{Method: http.MethodPost, Path: "/sys_apis/update_ignore", Group: "系统模块-API", Description: "修改忽略API"},
			// 系统模块-菜单
			{Method: http.MethodGet, Path: "/sys_menu/get", Group: "系统模块-菜单", Description: "获取菜单"},
			{Method: http.MethodGet, Path: "/sys_menu/list", Group: "系统模块-菜单", Description: "菜单列表"},
			{Method: http.MethodPost, Path: "/sys_menu/create", Group: "系统模块-菜单", Description: "创建菜单"},
			{Method: http.MethodPost, Path: "/sys_menu/delete", Group: "系统模块-菜单", Description: "删除菜单"},
			{Method: http.MethodPost, Path: "/sys_menu/edit", Group: "系统模块-菜单", Description: "修改菜单"},
			// 系统模块-角色
			{Method: http.MethodPost, Path: "/sys_role/permission", Group: "系统模块-角色", Description: "权限配置"},
			{Method: http.MethodGet, Path: "/sys_role/list", Group: "系统模块-角色", Description: "角色列表"},
			{Method: http.MethodGet, Path: "/sys_role/get", Group: "系统模块-角色", Description: "获取角色"},
			{Method: http.MethodPost, Path: "/sys_role/create", Group: "系统模块-角色", Description: "创建角色"},
			{Method: http.MethodPost, Path: "/sys_role/delete", Group: "系统模块-角色", Description: "删除角色"},
			{Method: http.MethodPost, Path: "/sys_role/edit", Group: "系统模块-角色", Description: "修改角色"},
			// 系统模块-定时任务
			{Method: http.MethodGet, Path: "/sys_cron_job/get", Group: "系统模块-定时任务", Description: "获取定时任务"},
			{Method: http.MethodGet, Path: "/sys_cron_job/get_handlers", Group: "系统模块-定时任务", Description: "可用任务列表"},
			{Method: http.MethodPost, Path: "/sys_cron_job/list", Group: "系统模块-定时任务", Description: "定时任务列表"},
			{Method: http.MethodPost, Path: "/sys_cron_job/create", Group: "系统模块-定时任务", Description: "创建定时任务"},
			{Method: http.MethodPost, Path: "/sys_cron_job/delete", Group: "系统模块-定时任务", Description: "删除定时任务"},
			{Method: http.MethodPost, Path: "/sys_cron_job/edit", Group: "系统模块-定时任务", Description: "修改定时任务"},
			{Method: http.MethodPost, Path: "/sys_cron_job/logs", Group: "系统模块-定时任务", Description: "任务日志列表"},
			//系统模块-开发工具
			{Method: http.MethodGet, Path: "/sys_devtools/history", Group: "系统模块-开发工具", Description: "代码生成历史"},
			{Method: http.MethodPost, Path: "/sys_devtools/generate", Group: "系统模块-开发工具", Description: "生成代码"},
			{Method: http.MethodPost, Path: "/sys_devtools/preview", Group: "系统模块-开发工具", Description: "预览代码"},
			{Method: http.MethodPost, Path: "/sys_devtools/delete_history", Group: "系统模块-开发工具", Description: "删除生成历史"},
		}
		return global.DB.Create(defaultApis).Error
	}
	ignoreCount := int64(0)
	if err := global.DB.Model(SysIgnoreApi{}).Count(&ignoreCount).Error; err != nil {
		return err
	}
	if ignoreCount == 0 {
		var defaultIgnoreApis = []SysIgnoreApi{
			{Method: http.MethodPost, Path: "/sys_auth/login"},
		}
		return global.DB.Create(defaultIgnoreApis).Error
	}
	return nil
}
