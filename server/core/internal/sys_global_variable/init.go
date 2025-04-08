package sys_global_variable

import (
	"gitee.com/mxcker/go-admin/server/core/global"
)

func InitData() error {
	var list []*SysGlobalVariable
	if err := global.DB.Model(&SysGlobalVariable{}).Find(&list).Error; err != nil {
		return err
	}
	for _, v := range list {
		global.Vars.Set(v.Key, v.Value)
	}
	return nil
}
