package modules

import (
	"gitee.com/mxcker/go-admin/server/core"
	"gitee.com/mxcker/go-admin/server/modules/test"
)

func Init() {
	core.AddModule("test", new(test.Mounter))
}
