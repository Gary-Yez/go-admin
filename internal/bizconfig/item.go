package bizconfig

import (
	"encoding/json"
	"errors"
	"reflect"
)

type Value interface {
	string | bool | int | float64 | []string
}

// Item 绑定单项配置的 Key，读取时才访问缓存或数据库。
type Item[T Value] struct {
	key     string
	initial T
}

func (item *Item[T]) Get() (T, error) {
	var value T
	if item.key == "" {
		return value, errors.New("配置项尚未注册")
	}
	raw, err := ReadValue(item.key)
	if err != nil {
		return value, err
	}
	if err := decodeStoredValue(reflect.TypeFor[T]().String(), raw, &value); err != nil {
		return value, err
	}
	return value, nil
}

func (item *Item[T]) MustGet() T {
	value, err := item.Get()
	if err != nil {
		panic(err)
	}
	return value
}

// InitialValue 供配置初始化构造动态初始值，不绑定 Key 或修改存储。
func InitialValue[T Value](value T) Item[T] { return Item[T]{initial: value} }

func (item *Item[T]) bind(key string, raw json.RawMessage) error {
	item.key = key
	return decodeStoredValue(reflect.TypeFor[T]().String(), raw, &item.initial)
}

func (item *Item[T]) valueType() reflect.Type               { return reflect.TypeFor[T]() }
func (item *Item[T]) defaultJSON() (json.RawMessage, error) { return json.Marshal(item.initial) }

type configItem interface {
	bind(string, json.RawMessage) error
	valueType() reflect.Type
	defaultJSON() (json.RawMessage, error)
}
