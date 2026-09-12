package bizconfig

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Gary-Yez/go-admin/internal/state"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ValidateValue(kind string, value json.RawMessage) error {
	var target any
	switch kind {
	case "string":
		target = new(string)
	case "bool":
		target = new(bool)
	case "int":
		target = new(int64)
	case "float64":
		target = new(float64)
	case "[]string":
		target = new([]string)
	default:
		return fmt.Errorf("不支持的配置类型：%s", kind)
	}
	return decodeStoredValue(kind, value, target)
}

// 直接解码到目标字段，同时校验空值、类型及整数范围。
func decodeStoredValue(kind string, value json.RawMessage, target any) error {
	invalid := fmt.Errorf("配置值必须为 %s，且不能为 null", kind)
	if strings.TrimSpace(string(value)) == "null" {
		return invalid
	}
	if err := json.Unmarshal(value, target); err != nil {
		return invalid
	}
	var number int64
	switch target := target.(type) {
	case *int:
		number = int64(*target)
	case *int64:
		number = *target
	default:
		return nil
	}
	if number > 9007199254740991 || number < -9007199254740991 {
		return errors.New("整数超出安全范围")
	}
	return nil
}

// Definitions 读取数据库中的配置定义，供管理页面及生成器使用。
func Definitions() ([]Definition, error) {
	var rows []SysConfigValue
	if err := state.DB().Select("config_key", "config_group", "label", "description", "type", "default_value", "sort").Order("sort, config_key").Find(&rows).Error; err != nil {
		return nil, err
	}
	definitions := make([]Definition, 0, len(rows))
	for _, row := range rows {
		definitions = append(definitions, row.definition())
	}
	return definitions, nil
}

func definitionRow(definition Definition) SysConfigValue {
	return SysConfigValue{Sort: definition.Sort, Key: definition.Key, Group: definition.Group, Label: definition.Label, Description: definition.Description, Type: definition.Type, DefaultValue: string(definition.Default), Value: string(definition.Default)}
}

// 代码定义仅用于首次部署补建。已有定义和值不随实例版本变化而覆盖。
func seedDefinition(definition Definition) error {
	if strings.TrimSpace(definition.Key) == "" || len(definition.Key) > 191 {
		return errors.New("配置标识不能为空且不能超过 191 字节")
	}
	if err := ValidateValue(definition.Type, definition.Default); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	lock, err := lockValue(ctx, definition.Key)
	if err != nil {
		return err
	}
	defer lock.Unlock()
	row := definitionRow(definition)
	if err := state.DB().WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&row).Error; err != nil {
		return err
	}
	if err := state.DB().WithContext(ctx).First(&row, "config_key = ?", definition.Key).Error; err != nil {
		return err
	}
	if row.Type != definition.Type {
		return fmt.Errorf("配置 %s 的代码类型 %s 与数据库类型 %s 不一致，请使用新键迁移", row.Key, definition.Type, row.Type)
	}
	if err := ValidateValue(row.Type, json.RawMessage(row.Value)); err != nil {
		return err
	}
	if err := validateConfigValue(row.Key, json.RawMessage(row.Value)); err != nil {
		return err
	}
	// 顺序随当前代码声明同步，其他定义和实际值保持不变。
	if row.Sort != definition.Sort {
		if err := state.DB().WithContext(ctx).Model(&row).UpdateColumn("sort", definition.Sort).Error; err != nil {
			return err
		}
	}
	return cacheValue(row.Key, json.RawMessage(row.Value))
}

// SaveDefinitions 在同一个事务内保存定义并执行文件写入，失败时回滚数据库。
func SaveDefinitions(definitions []Definition, apply func() error) error {
	return state.DB().Transaction(func(tx *gorm.DB) error {
		if err := saveDefinitions(tx, definitions); err != nil {
			return err
		}
		return apply()
	})
}

// 修改定义保留实际值，不删除未出现在本地代码中的配置。
func saveDefinitions(tx *gorm.DB, definitions []Definition) error {
	for _, field := range definitions {
		key := field.Key
		var row SysConfigValue
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&row, "config_key = ?", key).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			row = definitionRow(field)
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
		if row.Type != field.Type {
			return fmt.Errorf("配置 %s 不允许修改已有类型", key)
		}
		if row.Sort == field.Sort && row.Group == field.Group && row.Label == field.Label && row.Description == field.Description && row.DefaultValue == string(field.Default) {
			continue
		}
		result := tx.Model(&SysConfigValue{}).Where("config_key = ?", key).
			Updates(map[string]any{"sort": field.Sort, "config_group": field.Group, "label": field.Label, "description": field.Description, "default_value": string(field.Default)})
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
