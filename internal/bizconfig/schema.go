package bizconfig

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type SchemaField struct {
	Index []int
	Definition
}

type Schema struct {
	Type   reflect.Type
	Fields []SchemaField
}

// ParseSchema 只在启动注册时反射结构，读取时复用字段索引。
func ParseSchema(kind reflect.Type) (*Schema, error) {
	if kind.Kind() != reflect.Struct {
		return nil, fmt.Errorf("配置类型必须为结构体")
	}
	schema := &Schema{Type: kind}
	keys, names := map[string]bool{}, map[string]bool{}
	var visit func(reflect.Type, []int) error
	visit = func(kind reflect.Type, parent []int) error {
		for i := 0; i < kind.NumField(); i++ {
			field := kind.Field(i)
			key := field.Tag.Get("config")
			if key == "-" {
				continue
			}
			if !field.IsExported() {
				return fmt.Errorf("配置字段 %s 必须导出", field.Name)
			}
			index := append(append([]int(nil), parent...), i)
			if field.Anonymous && field.Type.Kind() == reflect.Struct && key == "" {
				if err := visit(field.Type, index); err != nil {
					return err
				}
				continue
			}
			if strings.TrimSpace(key) == "" || len(key) > 191 {
				return fmt.Errorf("字段 %s 需要有效的 config 标识", field.Name)
			}
			if keys[key] || names[field.Name] {
				return fmt.Errorf("配置字段或标识重复：%s", field.Name)
			}
			keys[key], names[field.Name] = true, true
			item, ok := reflect.New(field.Type).Interface().(configItem)
			if !ok {
				return fmt.Errorf("配置字段 %s 请使用 admin.ConfigItem[T]", field.Name)
			}
			name := item.valueType().String()
			value, present := field.Tag.Lookup("default")
			var raw json.RawMessage
			if name == "string" {
				raw, _ = json.Marshal(value)
			} else {
				if !present {
					switch name {
					case "bool":
						value = "false"
					case "int", "float64":
						value = "0"
					case "[]string":
						value = "[]"
					}
				}
				raw = json.RawMessage(value)
			}
			if err := ValidateValue(name, raw); err != nil {
				return fmt.Errorf("字段 %s：%w", field.Name, err)
			}
			label := field.Tag.Get("label")
			if label == "" {
				label = field.Name
			}
			schema.Fields = append(schema.Fields, SchemaField{Index: index, Definition: Definition{Key: key, Group: strings.TrimSpace(field.Tag.Get("group")), Label: label, Description: field.Tag.Get("description"), Type: name, Default: raw}})
		}
		return nil
	}
	if err := visit(kind, nil); err != nil {
		return nil, err
	}
	for i := range schema.Fields {
		schema.Fields[i].Sort = i
	}
	return schema, nil
}

func (schema *Schema) Definitions() []Definition {
	list := make([]Definition, 0, len(schema.Fields))
	for _, field := range schema.Fields {
		list = append(list, field.Definition)
	}
	return list
}

// initialDefinitions 按标签默认值、业务 Init 的顺序生成初始值。
func (schema *Schema) initialDefinitions() ([]Definition, error) {
	value := reflect.New(schema.Type)
	if err := schema.Bind(value.Elem()); err != nil {
		return nil, err
	}

	if initializer, ok := value.Interface().(interface{ Init() error }); ok {
		if err := initializer.Init(); err != nil {
			return nil, fmt.Errorf("初始化配置：%w", err)
		}
	}
	definitions := schema.Definitions()
	for i, field := range schema.Fields {
		item := value.Elem().FieldByIndex(field.Index).Addr().Interface().(configItem)
		raw, err := item.defaultJSON()
		if err != nil {
			return nil, fmt.Errorf("配置 %s：%w", field.Key, err)
		}
		if err := ValidateValue(field.Type, raw); err != nil {
			return nil, fmt.Errorf("配置 %s：%w", field.Key, err)
		}
		definitions[i].Default = raw
	}
	return definitions, nil
}

// Bind 只绑定 Key 和默认值，不读取数据库或缓存。
func (schema *Schema) Bind(target reflect.Value) error {
	for _, field := range schema.Fields {
		item := target.FieldByIndex(field.Index).Addr().Interface().(configItem)
		if err := item.bind(field.Key, field.Default); err != nil {
			return fmt.Errorf("配置 %s：%w", field.Key, err)
		}
	}
	return nil
}
