package cacher

import (
	"encoding/json"
	"errors"
	"github.com/patrickmn/go-cache"
	"time"
)

type memoryLock struct {
	key        string
	cache      *cache.Cache
	done       chan struct{}
	expireTime time.Duration
}

func (l *memoryLock) autoExtend(extendInterval time.Duration) {
	ticker := time.NewTicker(extendInterval)
	defer ticker.Stop()
	for {
		select {
		case <-l.done:
			return
		case <-ticker.C:
			err := l.cache.Replace(l.key, time.Now(), l.expireTime)
			if err != nil {
				return
			}
		}
	}
}

func (l *memoryLock) Unlock() error {
	close(l.done)
	l.cache.Delete(l.key)
	return nil
}
func NewMemoryCache(cleanupInterval time.Duration) (Cache, error) {
	return &memoryCache{
		cache: cache.New(cache.NoExpiration, cleanupInterval),
	}, nil
}

type memoryCache struct {
	cache *cache.Cache
}

func (m *memoryCache) Client() interface{} {
	return m.cache
}

func (m *memoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	m.cache.Set(key, value, ttl)
	return nil
}

func (m *memoryCache) Get(key string) (value string, err error) {
	v, ok := m.cache.Get(key)
	if !ok {
		return "", ErrCacheNotFound
	}
	value, ok = v.(string)
	if !ok {
		return "", errors.New("value is not string")
	}
	return value, nil
}

func (m *memoryCache) MustGet(key string) string {
	v, ok := m.cache.Get(key)
	if !ok {
		return ""
	}
	value, ok := v.(string)
	if !ok {
		return ""
	}
	return value
}

func (m *memoryCache) GetBool(key string) (value bool, err error) {
	v, ok := m.cache.Get(key)
	if !ok {
		return false, ErrCacheNotFound
	}
	value, ok = v.(bool)
	if !ok {
		return false, errors.New("value is not bool")
	}
	return value, nil
}

func (m *memoryCache) MustGetBool(key string) bool {
	v, ok := m.cache.Get(key)
	if !ok {
		return false
	}
	value, ok := v.(bool)
	if !ok {
		return false
	}
	return value
}

func (m *memoryCache) GetInt(key string) (value int, err error) {
	v, ok := m.cache.Get(key)
	if !ok {
		return 0, ErrCacheNotFound
	}
	value, ok = v.(int)
	if !ok {
		return 0, errors.New("value is not int")
	}
	return value, nil
}

func (m *memoryCache) MustGetInt(key string) int {
	v, ok := m.cache.Get(key)
	if !ok {
		return 0
	}
	value, ok := v.(int)
	if !ok {
		return 0
	}
	return value
}

func (m *memoryCache) Del(key string) error {
	m.cache.Delete(key)
	return nil
}

func (m *memoryCache) Exists(key string) (bool, error) {
	_, ok := m.cache.Get(key)
	return ok, nil
}

func (m *memoryCache) SetJSON(key string, data any, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	m.cache.Set(key, string(jsonData), ttl)
	return nil
}

func (m *memoryCache) GetJSON(key string, dest any) error {
	v, ok := m.cache.Get(key)
	if !ok {
		return ErrCacheNotFound
	}
	str, ok := v.(string)
	if !ok {
		return errors.New("value is not string")
	}
	return json.Unmarshal([]byte(str), dest)
}

func (m *memoryCache) Lock(key string, expiration time.Duration, extendInterval time.Duration) (Lock, error) {
	err := m.cache.Add(key, time.Now(), expiration)
	if err != nil {
		return nil, err
	}
	l := &memoryLock{
		cache:      m.cache,
		key:        key,
		done:       make(chan struct{}),
		expireTime: expiration,
	}
	if extendInterval > 0 {
		go l.autoExtend(extendInterval)
	}
	return l, nil
}
