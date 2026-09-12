package sys_apis

import (
	"github.com/Gary-Yez/go-admin/internal/state"
	"net/http"
)

func InitData() error {
	count := int64(0)
	if err := state.DB().Model(SysApi{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {

		if err := state.DB().Create(defaultAPIs()).Error; err != nil {
			return err
		}
	}
	return nil
}

// defaultAPIs 同时用于首次填充和启动时补录新接口的分组、说明。
func defaultAPIs() []SysApi {
	return []SysApi{
		{Method: http.MethodGet, Path: "/sys_config/values", Group: "系统运维-配置管理", Description: "配置项列表"},
		{Method: http.MethodPost, Path: "/sys_config/update_value", Group: "系统运维-配置管理", Description: "修改配置项"},
		{Method: http.MethodPost, Path: "/sys_config/reset_value", Group: "系统运维-配置管理", Description: "重置配置项"},
		{Method: http.MethodPost, Path: "/sys_config/sync_cache", Group: "系统运维-配置管理", Description: "同步配置缓存"},
		{Method: http.MethodPost, Path: "/sys_config/cleanup_invalid", Group: "系统运维-配置管理", Description: "删除无效配置"},
		{Method: http.MethodGet, Path: "/sys_monitor/list", Group: "系统运维-节点监控", Description: "节点资源与 Go Runtime 监控"},
		{Method: http.MethodPost, Path: "/sys_api_token/list", Group: "权限管理-API密钥管理", Description: "全局 API 密钥列表"},
		{Method: http.MethodPost, Path: "/sys_api_token/delete", Group: "权限管理-API密钥管理", Description: "删除全局 API 密钥"},
		{Method: http.MethodPost, Path: "/sys_login_log/list", Group: "系统运维-登录日志", Description: "登录日志列表"},
		{Method: http.MethodPost, Path: "/sys_login_log/delete", Group: "系统运维-登录日志", Description: "删除登录日志"},
		{Method: http.MethodPost, Path: "/sys_login_log/cleanup", Group: "系统运维-登录日志", Description: "清理历史登录日志"},
		// 权限管理-管理员
		{Method: http.MethodPost, Path: "/sys_admin/create", Group: "权限管理-管理员", Description: "创建管理员"},
		{Method: http.MethodPost, Path: "/sys_admin/delete", Group: "权限管理-管理员", Description: "删除管理员"},
		{Method: http.MethodPost, Path: "/sys_admin/edit", Group: "权限管理-管理员", Description: "修改管理员"},
		{Method: http.MethodGet, Path: "/sys_admin/list", Group: "权限管理-管理员", Description: "管理员列表"},
		// 权限管理-API
		{Method: http.MethodPost, Path: "/sys_apis/delete", Group: "权限管理-API", Description: "删除API"},
		{Method: http.MethodPost, Path: "/sys_apis/edit", Group: "权限管理-API", Description: "修改API"},
		{Method: http.MethodGet, Path: "/sys_apis/get_groups", Group: "权限管理-API", Description: "API分组列表"},
		{Method: http.MethodPost, Path: "/sys_apis/list", Group: "权限管理-API", Description: "API列表"},
		{Method: http.MethodGet, Path: "/sys_apis/invalid_apis", Group: "权限管理-API", Description: "查询失效API"},
		// 权限管理-菜单
		{Method: http.MethodGet, Path: "/sys_menu/list", Group: "权限管理-菜单", Description: "菜单列表"},
		{Method: http.MethodPost, Path: "/sys_menu/create", Group: "权限管理-菜单", Description: "创建菜单"},
		{Method: http.MethodPost, Path: "/sys_menu/delete", Group: "权限管理-菜单", Description: "删除菜单"},
		{Method: http.MethodPost, Path: "/sys_menu/edit", Group: "权限管理-菜单", Description: "修改菜单"},
		{Method: http.MethodPost, Path: "/sys_menu/sort", Group: "权限管理-菜单", Description: "菜单拖拽排序"},
		// 权限管理-角色
		{Method: http.MethodPost, Path: "/sys_role/permission", Group: "权限管理-角色", Description: "权限配置"},
		{Method: http.MethodGet, Path: "/sys_role/list", Group: "权限管理-角色", Description: "角色列表"},
		{Method: http.MethodGet, Path: "/sys_role/get", Group: "权限管理-角色", Description: "获取角色"},
		{Method: http.MethodPost, Path: "/sys_role/create", Group: "权限管理-角色", Description: "创建角色"},
		{Method: http.MethodPost, Path: "/sys_role/copy", Group: "权限管理-角色", Description: "拷贝角色及权限"},
		{Method: http.MethodPost, Path: "/sys_role/delete", Group: "权限管理-角色", Description: "删除角色"},
		{Method: http.MethodPost, Path: "/sys_role/edit", Group: "权限管理-角色", Description: "修改角色"},
		// 系统运维-计划任务
		{Method: http.MethodGet, Path: "/sys_cron_job/get_handlers", Group: "系统运维-计划任务", Description: "可用任务列表"},
		{Method: http.MethodPost, Path: "/sys_cron_job/list", Group: "系统运维-计划任务", Description: "定时任务列表"},
		{Method: http.MethodPost, Path: "/sys_cron_job/create", Group: "系统运维-计划任务", Description: "创建定时任务"},
		{Method: http.MethodPost, Path: "/sys_cron_job/delete", Group: "系统运维-计划任务", Description: "删除定时任务"},
		{Method: http.MethodPost, Path: "/sys_cron_job/edit", Group: "系统运维-计划任务", Description: "修改定时任务"},
		{Method: http.MethodPost, Path: "/sys_cron_job/logs", Group: "系统运维-计划任务", Description: "任务日志列表"},
	}
}
