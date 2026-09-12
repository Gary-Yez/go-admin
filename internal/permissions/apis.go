package permissions

import "net/http"

type API struct {
	Method string
	Path   string
}

// AuthenticatedAPIs 是操作当前用户自身的基础接口，不参与角色授权。
// 精确匹配方法和路径，新增鉴权模块接口不会自动获得豁免。
func AuthenticatedAPIs() []API {
	return []API{
		{Method: http.MethodGet, Path: "/sys_auth/me"},
		{Method: http.MethodGet, Path: "/sys_auth/password_policy"},
		{Method: http.MethodPost, Path: "/sys_auth/switch_role"},
		{Method: http.MethodPost, Path: "/sys_auth/change_info"},
		{Method: http.MethodPost, Path: "/sys_auth/change_password"},
		{Method: http.MethodGet, Path: "/sys_api_token/mine/list"},
		{Method: http.MethodPost, Path: "/sys_api_token/mine/save"},
		{Method: http.MethodPost, Path: "/sys_api_token/mine/delete"},
	}
}

func IsAuthenticatedAPI(method, path string) bool {
	for _, api := range AuthenticatedAPIs() {
		if api.Method == method && api.Path == path {
			return true
		}
	}
	return false
}
