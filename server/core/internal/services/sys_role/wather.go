package sys_role

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisWatcher struct {
	cb      func(string)
	rdb     *redis.Client
	channel string
}

func NewRedisWatcher(rdb *redis.Client, channel string) *RedisWatcher {
	return &RedisWatcher{rdb: rdb, channel: channel}
}

func (w *RedisWatcher) SetUpdateCallback(fn func(string)) error {
	w.cb = fn
	go w.subscribe()
	return nil
}

func (w *RedisWatcher) Update() error {
	// 发布消息通知其它实例
	return w.rdb.Publish(context.Background(), w.channel, "update").Err()
}

func (w *RedisWatcher) subscribe() {
	sub := w.rdb.Subscribe(context.Background(), w.channel)
	for msg := range sub.Channel() {
		if w.cb != nil {
			w.cb(msg.Payload)
		}
	}
}
