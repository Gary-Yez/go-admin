package bizconfig

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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

func ValidateDefinition(definition Definition) error {
	if err := ValidateValue(definition.Type, definition.Default); err != nil {
		return err
	}
	requiresOptions := false
	switch definition.Control {
	case "":
	case "textarea", "password":
		if definition.Type != "string" {
			return fmt.Errorf("%s 控件只支持 string 类型", definition.Control)
		}
	case "select":
		if definition.Type != "string" && definition.Type != "int" && definition.Type != "float64" {
			return fmt.Errorf("select 控件不支持 %s 类型", definition.Type)
		}
		requiresOptions = true
	case "radio":
		if definition.Type != "string" && definition.Type != "int" && definition.Type != "float64" {
			return fmt.Errorf("radio 控件不支持 %s 类型", definition.Type)
		}
		requiresOptions = true
	case "multi-select":
		if definition.Type != "[]string" {
			return errors.New("multi-select 控件只支持 []string 类型")
		}
		requiresOptions = true
	default:
		return fmt.Errorf("不支持的配置控件：%s", definition.Control)
	}
	if definition.Control == "textarea" {
		if definition.Rows < 1 || definition.Rows > 20 {
			return errors.New("多行文本框默认行数必须在 1 到 20 之间")
		}
	} else if definition.Rows != 0 {
		return errors.New("只有多行文本框可以设置默认行数")
	}
	if !requiresOptions {
		if len(definition.Options) != 0 {
			return errors.New("当前配置控件不能设置候选项")
		}
		return nil
	}
	if len(definition.Options) == 0 {
		return errors.New("选择控件至少需要一个候选项")
	}
	optionType := definition.Type
	if definition.Control == "multi-select" {
		optionType = "string"
	}
	seen := map[string]bool{}
	for _, option := range definition.Options {
		if strings.TrimSpace(option.Label) == "" {
			return errors.New("下拉候选项名称不能为空")
		}
		if err := ValidateValue(optionType, option.Value); err != nil {
			return fmt.Errorf("候选项 %s：%w", option.Label, err)
		}
		var decoded any
		_ = json.Unmarshal(option.Value, &decoded)
		canonical, _ := json.Marshal(decoded)
		key := string(canonical)
		if seen[key] {
			return errors.New("下拉候选值不能重复")
		}
		seen[key] = true
	}
	return ValidateOptionValue(definition, definition.Default)
}

func ValidateOptionValue(definition Definition, value json.RawMessage) error {
	if definition.Control != "select" && definition.Control != "radio" && definition.Control != "multi-select" {
		return nil
	}
	var actual any
	if err := json.Unmarshal(value, &actual); err != nil {
		return err
	}
	if definition.Control == "multi-select" {
		values, ok := actual.([]any)
		if !ok {
			return errors.New("多选配置值必须为字符串列表")
		}
		for _, value := range values {
			if !containsOptionValue(definition.Options, value) {
				return errors.New("配置值包含不在候选项中的内容")
			}
		}
		return nil
	}
	if containsOptionValue(definition.Options, actual) {
		return nil
	}
	return errors.New("配置值不在候选项中")
}

func containsOptionValue(options []Option, actual any) bool {
	for _, option := range options {
		var candidate any
		if json.Unmarshal(option.Value, &candidate) == nil && reflect.DeepEqual(actual, candidate) {
			return true
		}
	}
	return false
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
	if err := state.DB().Select("config_key", "config_group", "label", "description", "type", "control", "rows", "config_options", "default_value", "sort").Order("sort, config_key").Find(&rows).Error; err != nil {
		return nil, err
	}
	definitions := make([]Definition, 0, len(rows))
	for _, row := range rows {
		definitions = append(definitions, row.definition())
	}
	return definitions, nil
}

func definitionRow(definition Definition) SysConfigValue {
	options, _ := json.Marshal(definition.Options)
	return SysConfigValue{Sort: definition.Sort, Key: definition.Key, Group: definition.Group, Label: definition.Label, Description: definition.Description, Type: definition.Type, Control: definition.Control, Rows: definition.Rows, OptionsValue: string(options), DefaultValue: string(definition.Default), Value: string(definition.Default)}
}

// 代码定义仅用于首次部署补建。已有定义和值不随实例版本变化而覆盖。
func seedDefinition(definition Definition) error {
	if strings.TrimSpace(definition.Key) == "" || len(definition.Key) > 191 {
		return errors.New("配置标识不能为空且不能超过 191 字节")
	}
	if err := ValidateDefinition(definition); err != nil {
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
		if err := ValidateDefinition(field); err != nil {
			return fmt.Errorf("配置 %s：%w", key, err)
		}
		if err := ValidateOptionValue(field, json.RawMessage(row.Value)); err != nil {
			return fmt.Errorf("配置 %s 的当前值：%w", key, err)
		}
		options, _ := json.Marshal(field.Options)
		if row.Sort == field.Sort && row.Group == field.Group && row.Label == field.Label && row.Description == field.Description && row.Control == field.Control && row.Rows == field.Rows && row.OptionsValue == string(options) && row.DefaultValue == string(field.Default) {
			continue
		}
		result := tx.Model(&SysConfigValue{}).Where("config_key = ?", key).
			Updates(map[string]any{"sort": field.Sort, "config_group": field.Group, "label": field.Label, "description": field.Description, "control": field.Control, "rows": field.Rows, "config_options": string(options), "default_value": string(field.Default)})
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
