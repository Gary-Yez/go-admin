package bizconfig

import (
	"context"
	"encoding/json"

	"fmt"

	"github.com/Gary-Yez/go-admin/internal/state"
)

// 启动成功后只读，用于判断当前实例仍在使用的配置。
var registeredKeys map[string]struct{}

// 业务校验由调用方提供，运行期间保持不变。
var validateConfigValue = func(string, json.RawMessage) error { return nil }

func Initialize(schema *Schema, validate func(string, json.RawMessage) error) error {
	if validate != nil {
		validateConfigValue = validate
	}
	definitions, err := schema.initialDefinitions()
	if err != nil {
		return err
	}
	if err := state.DB().AutoMigrate(&SysConfigValue{}); err != nil {
		return err
	}
	keys := make(map[string]struct{}, len(definitions))
	for _, definition := range definitions {
		if err := seedDefinition(definition); err != nil {
			return err
		}
		keys[definition.Key] = struct{}{}
	}
	registeredKeys = keys
	return nil
}

func Values(ctx context.Context) ([]ValueRow, error) {
	var rows []SysConfigValue
	if err := state.DB().WithContext(ctx).Order("sort, config_key").Find(&rows).Error; err != nil {
		return nil, err
	}
	list := make([]ValueRow, 0, len(rows))
	for _, row := range rows {
		value := json.RawMessage(row.Value)
		if !json.Valid(value) {
			return nil, fmt.Errorf("配置 %s 的存储值不是有效 JSON", row.Key)
		}
		definition := row.definition()
		if !json.Valid(definition.Default) {
			return nil, fmt.Errorf("配置 %s 的默认值不是有效 JSON", row.Key)
		}
		_, registered := registeredKeys[row.Key]
		list = append(list, ValueRow{Definition: definition, Value: value, UpdatedAt: row.UpdatedAt, Registered: registered})
	}
	return list, nil
}

func UpdateValue(ctx context.Context, key string, value json.RawMessage, reset bool) error {
	lock, err := lockValue(ctx, key)
	if err != nil {
		return err
	}
	defer lock.Unlock()
	var row SysConfigValue
	if err := state.DB().WithContext(ctx).First(&row, "config_key = ?", key).Error; err != nil {
		return err
	}
	if reset {
		value = json.RawMessage(row.DefaultValue)
	}
	if err := ValidateValue(row.Type, value); err != nil {
		return err
	}
	if err := ValidateOptionValue(row.definition(), value); err != nil {
		return err
	}
	if err := validateConfigValue(row.Key, value); err != nil {
		return err
	}
	result := state.DB().WithContext(ctx).Model(&SysConfigValue{}).
		Where("config_key = ?", key).
		Updates(map[string]any{"value": string(value)})
	if result.Error != nil {
		return result.Error
	}
	if err := cacheValue(key, value); err != nil {
		return fmt.Errorf("配置已保存到数据库，但缓存更新失败：%w", err)
	}
	return nil
}
