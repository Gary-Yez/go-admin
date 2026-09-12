package sys_menu

import (
	"fmt"
	"github.com/Gary-Yez/go-admin/internal/state"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Definition struct {
	Key, Name, ParentKey, Icon, Path, Component string
	Sort                                        int
	Hidden                                      bool
}

// Sync 只补充缺失菜单，不覆盖后台调整；所有父子关系在事务中建立。
func Sync(definitions []Definition) error {
	return state.DB().Transaction(func(tx *gorm.DB) error {
		pending := append([]Definition(nil), definitions...)
		for len(pending) > 0 {
			next := make([]Definition, 0)
			for _, item := range pending {
				if item.Key == "" || item.Name == "" || item.Key == item.ParentKey {
					return fmt.Errorf("菜单配置无效：%s", item.Key)
				}
				var existing SysMenu
				if err := tx.Where(map[string]interface{}{"key": item.Key}).Limit(1).Find(&existing).Error; err != nil {
					return err
				}
				if existing.Id != 0 {
					continue
				}
				menu := SysMenu{Key: item.Key, Name: item.Name, Icon: item.Icon, Path: item.Path, Component: item.Component, Sort: item.Sort, Hidden: item.Hidden}
				if item.ParentKey != "" {
					var parent SysMenu
					if err := tx.Where(map[string]interface{}{"key": item.ParentKey}).Limit(1).Find(&parent).Error; err != nil {
						return err
					}
					if parent.Id == 0 {
						next = append(next, item)
						continue
					}
					menu.ParentId = &parent.Id
				}
				if err := tx.Omit(clause.Associations).Clauses(clause.OnConflict{DoNothing: true}).Create(&menu).Error; err != nil {
					return err
				}
			}
			if len(next) == len(pending) {
				return fmt.Errorf("菜单父级不存在或存在循环：%s", next[0].Key)
			}
			pending = next
		}
		return nil
	})
}
