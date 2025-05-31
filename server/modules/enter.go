package modules

import (
	"gitee.com/mxcker/go-admin/server/modules/test"
	"gitee.com/mxcker/go-admin/server/pkg/modular"
)

func Load(loader modular.Loader) error {
	loader.Mount("test", new(test.Mounter))
	return nil
}
