package vars_cache

import (
	"strconv"
	"sync"
)

type memoryCache struct {
	data sync.Map
}

func (m *memoryCache) Set(key string, value string) {
	m.data.Store(key, value)
}

func (m *memoryCache) GetString(key string) string {
	v, ok := m.data.Load(key)
	if !ok {
		return ""
	}
	return v.(string)
}

func (m *memoryCache) GetInt(key string) int {
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

func (m *memoryCache) GetBool(key string) bool {
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
