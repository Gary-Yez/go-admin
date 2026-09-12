package sys_devtools

import (
	"fmt"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/system/sys_menu"
	"strings"
)

func prepareMenus(data *GenerateBody) error {
	data.menuDefinitions = nil
	if !data.CreateCURD || !data.CreateMenu {
		return nil
	}
	name := strings.TrimSpace(data.MenuName)
	if name == "" {
		name = data.ChineseModuleName
	}
	if len(name) > 100 || len(data.MenuIcon) > 100 {
		return fmt.Errorf("菜单名称或图标过长")
	}
	leaf := sys_menu.Definition{Key: data.ModuleName, Name: name, ParentKey: data.MenuParentKey, Icon: data.MenuIcon, Path: data.ModuleName, Component: "../views/" + data.ModuleName + "/index.vue"}
	var existing sys_menu.SysMenu
	if err := state.DB().Where("`key` = ?", leaf.Key).Limit(1).Find(&existing).Error; err != nil {
		return err
	}
	if existing.Id != 0 && existing.Component != leaf.Component {
		return fmt.Errorf("菜单标识 %s 已被其他页面使用", leaf.Key)
	}
	data.menuDefinitions = append(data.menuDefinitions, leaf)
	key := data.MenuParentKey
	seen := map[string]bool{leaf.Key: true}
	for key != "" {
		if seen[key] {
			return fmt.Errorf("菜单父级存在循环")
		}
		seen[key] = true
		var parent sys_menu.SysMenu
		if err := state.DB().Where("`key` = ?", key).First(&parent).Error; err != nil {
			return fmt.Errorf("读取父菜单 %s：%w", key, err)
		}
		definition := sys_menu.Definition{Key: parent.Key, Name: parent.Name, Icon: parent.Icon, Path: parent.Path, Component: parent.Component, Sort: parent.Sort, Hidden: parent.Hidden}
		key = ""
		if parent.ParentId != nil {
			var ancestor sys_menu.SysMenu
			if err := state.DB().First(&ancestor, *parent.ParentId).Error; err != nil {
				return err
			}
			if ancestor.Key == "" {
				return fmt.Errorf("父菜单 %s 的上级缺少唯一标识，请先完善菜单配置", parent.Name)
			}
			key = ancestor.Key
		}
		definition.ParentKey = key
		data.menuDefinitions = append(data.menuDefinitions, definition)
	}
	return nil
}
