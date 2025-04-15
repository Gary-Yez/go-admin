package modules

import (
	"gitee.com/mxcker/go-admin/server/core/pkg/modular"
	"gitee.com/mxcker/go-admin/server/modules/test"
)

func Load(loader modular.Loader) error {
	loader.Add("test", new(test.Mounter))
	return nil
}
