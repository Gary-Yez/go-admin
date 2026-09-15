package initialization

import (
	"encoding/json"
	"reflect"

	"github.com/Gary-Yez/go-admin/internal/bizconfig"
	"github.com/Gary-Yez/go-admin/internal/system/sys_config"
	"github.com/Gary-Yez/go-admin/internal/system/sys_file"
	"github.com/Gary-Yez/go-admin/internal/utils"
)

// InitBusinessConfig 在数据库和缓存就绪后初始化业务配置，T 为框架内置配置类型。
func InitBusinessConfig[T any](schema *bizconfig.Schema) (*bizconfig.Schema, error) {
	base, err := bizconfig.ParseSchema(reflect.TypeFor[T]())
	if err != nil {
		return nil, err
	}
	if schema == nil {
		schema = base
	}
	if err := bizconfig.Initialize(schema, func(key string, value json.RawMessage) error {
		if err := utils.ValidateLoginValue(key, value); err != nil {
			return err
		}
		return sys_file.ValidateConfigValue(key, value)
	}); err != nil {
		return nil, err
	}
	sys_config.SetBaseSchema(base)
	return schema, nil
}
