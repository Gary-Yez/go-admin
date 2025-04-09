package modules

import (
	"gitee.com/mxcker/go-admin/server/core/types/module_loader"
	"gitee.com/mxcker/go-admin/server/modules/sys_task"
)

func Load(loader module_loader.Loader) error {
	loader.Add("sys_task", new(sys_task.Mounter))
	return nil
}
