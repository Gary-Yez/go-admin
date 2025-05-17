package cacher

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisLock struct {
	done chan struct{}
	mu   *redsync.Mutex
}

func (l *redisLock) autoExtend(extendInterval time.Duration) {
	ticker := time.NewTicker(extendInterval)
	defer ticker.Stop()
	for {
		select {
		case <-l.done:
			return
		case <-ticker.C:
			_, err := l.mu.Extend()
			if err != nil {
				return
			}
		}
	}
}

func (l *redisLock) Unlock() error {
	close(l.done)
	unlock, err := l.mu.Unlock()
	if err != nil || !unlock {
		return errors.New("redis lock unlock failed")
	}
	return nil
}

func NewRedisCache(options *redis.Options) (Cache, error) {
	rdb := &redisCache{
		ctx:    context.Background(),
		client: redis.NewClient(options),
	}
	pool := goredis.NewPool(rdb.client) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)
	rdb.lc = rs
	// 测试连通性
	if err := rdb.client.Ping(rdb.ctx).Err(); err != nil {
		return nil, errors.New("redis error:" + err.Error())
	}
	return rdb, nil
}

type redisCache struct {
	ctx    context.Context
	client *redis.Client
	lc     *redsync.Redsync
}

func (r *redisCache) Client() interface{} {
	return r.client
}

func (r *redisCache) Set(key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(r.ctx, key, value, ttl).Err()
}

func (r *redisCache) SetJSON(key string, data any, ttl time.Duration) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.Set(key, string(jsonBytes), ttl)
}

func (r *redisCache) Get(key string) (value string, err error) {
	value, err = r.client.Get(r.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return value, ErrCacheNotFound
	}
	return value, err
}

func (r *redisCache) MustGet(key string) string {
	value, _ := r.client.Get(r.ctx, key).Result()
	return value
}

func (r *redisCache) GetBool(key string) (value bool, err error) {
	value, err = r.client.Get(r.ctx, key).Bool()
	if errors.Is(err, redis.Nil) {
		return value, ErrCacheNotFound
	}
	return value, err
}

func (r *redisCache) MustGetBool(key string) bool {
	value, _ := r.client.Get(r.ctx, key).Bool()
	return value
}

func (r *redisCache) GetInt(key string) (value int, err error) {
	value, err = r.client.Get(r.ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		return value, ErrCacheNotFound
	}
	return value, err
}

func (r *redisCache) MustGetInt(key string) int {
	value, _ := r.client.Get(r.ctx, key).Int()
	return value
}

func (r *redisCache) GetJSON(key string, dest any) error {
	str, err := r.Get(key)
	if errors.Is(err, redis.Nil) {
		return ErrCacheNotFound
	} else if err != nil {
		return err
	}
	return json.Unmarshal([]byte(str), dest)
}

func (r *redisCache) Exists(key string) (bool, error) {
	res, err := r.client.Exists(r.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, ErrCacheNotFound
	}
	return res > 0, err
}

func (r *redisCache) Del(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *redisCache) Lock(key string, expiration time.Duration, extendInterval time.Duration) (Lock, error) {
	mu := r.lc.NewMutex(key, redsync.WithTries(1), redsync.WithExpiry(expiration))
	err := mu.LockContext(r.ctx)
	if err != nil {
		return nil, err
	}
	l := &redisLock{
		done: make(chan struct{}),
		mu:   mu,
	}
	if extendInterval > 0 {
		//需要自动续期
		go l.autoExtend(extendInterval)
	}
	return l, nil
}
