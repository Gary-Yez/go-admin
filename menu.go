package admin

import "github.com/Gary-Yez/go-admin/internal/system/sys_menu"

// MenuDefinition 是可随业务代码部署的菜单默认配置，父级使用稳定的 Key。
type MenuDefinition = sys_menu.Definition

// MenuProvider 是可选接口，不影响现有 Module 实现。
type MenuProvider interface {
	Menus() []MenuDefinition
}
