package sys_global_variable

import (
	"gitee.com/mxcker/go-admin/server/core/global"
	"gorm.io/gorm/clause"
)

func InitData() error {
	// 批量插入，冲突时忽略
	if err := global.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(Vars).Error; err != nil {
		return err
	}
	var list []*SysGlobalVariable
	if err := global.DB.Model(&SysGlobalVariable{}).Find(&list).Error; err != nil {
		return err
	}
	for _, v := range list {
		global.Vars.Set(v.Key, v.Value)
	}
	return nil
}
