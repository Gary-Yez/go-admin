package global

import (
	"strconv"
	"sync"
)

type variableMap struct {
	data sync.Map
}

func (m *variableMap) Set(key string, value string) {
	m.data.Store(key, value)
}

func (m *variableMap) GetString(key string) string {
	v, ok := m.data.Load(key)
	if !ok {
		return ""
	}
	return v.(string)
}

func (m *variableMap) GetInt(key string) int {
	v := m.GetString(key)
	if v == "" {
		return 0
	} else {
		result, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		return result
	}
}

func (m *variableMap) GetBool(key string) bool {
	v := m.GetString(key)
	if v == "" {
		return false
	} else {
		parseBool, err := strconv.ParseBool(v)
		if err != nil {
			return false
		}
		return parseBool
	}
}
