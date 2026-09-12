package monitor

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"sort"
	"strconv"
	"sync"
	"time"
)

// Store 是节点目录及快照存储，Redis 使用原子脚本，避免多节点争抢一个 JSON 或全库扫描。
type Store struct {
	client *redis.Client
	keys   []string
	mu     sync.RWMutex
	latest *Snapshot
}

func NewStore(client *redis.Client, namespace string) *Store {
	// 花括号确保目录和数据在 Redis Cluster 中位于同一 slot。
	return &Store{client: client, keys: []string{namespace + "cache:{monitor}:nodes", namespace + "cache:{monitor}:seen"}}
}

const pruneScript = `
local now=tonumber(redis.call('TIME')[1])
local expired=redis.call('ZRANGEBYSCORE',KEYS[2],'-inf',now-tonumber(ARGV[1]))
for _,id in ipairs(expired) do redis.call('HDEL',KEYS[1],id) end
redis.call('ZREMRANGEBYSCORE',KEYS[2],'-inf',now-tonumber(ARGV[1]))
`

var putScript = redis.NewScript(pruneScript + `
redis.call('HSET',KEYS[1],ARGV[2],cjson.encode({seen=now,data=ARGV[3]}))
redis.call('ZADD',KEYS[2],now,ARGV[2])
redis.call('EXPIRE',KEYS[1],ARGV[1])
redis.call('EXPIRE',KEYS[2],ARGV[1])
return 1
`)
var listScript = redis.NewScript(pruneScript + `
local rows=redis.call('HVALS',KEYS[1])
table.insert(rows,1,tostring(now))
return rows
`)

func (s *Store) Publish(ctx context.Context, row Snapshot) error {
	if s.client == nil {
		row.LastSeen = time.Now().UTC()
		s.mu.Lock()
		s.latest = &row
		s.mu.Unlock()
		return nil
	}
	data, err := json.Marshal(row)
	if err != nil {
		return err
	}
	return putScript.Run(ctx, s.client, s.keys, int(Retention.Seconds()), row.ID, string(data)).Err()
}
func (s *Store) List(ctx context.Context) ([]Snapshot, error) {
	rows := []Snapshot{}
	if s.client == nil {
		s.mu.RLock()
		if s.latest != nil {
			rows = append(rows, *s.latest)
		}
		s.mu.RUnlock()
		for i := range rows {
			rows[i].Online = time.Since(rows[i].LastSeen) <= OfflineAfter
		}
		return rows, nil
	}
	values, err := listScript.Run(ctx, s.client, s.keys, int(Retention.Seconds())).StringSlice()
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, errors.New("节点目录响应为空")
	}
	now, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return nil, err
	}
	for _, raw := range values[1:] {
		var entry struct {
			Seen int64  `json:"seen"`
			Data string `json:"data"`
		}
		if err := json.Unmarshal([]byte(raw), &entry); err != nil {
			return nil, err
		}
		var row Snapshot
		if err := json.Unmarshal([]byte(entry.Data), &row); err != nil {
			return nil, err
		}
		row.LastSeen = time.Unix(entry.Seen, 0).UTC()
		row.Online = now-entry.Seen <= int64(OfflineAfter.Seconds())
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Name == rows[j].Name {
			return rows[i].ID < rows[j].ID
		}
		return rows[i].Name < rows[j].Name
	})
	return rows, nil
}
