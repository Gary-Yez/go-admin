package sys_file

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Gary-Yez/go-admin/internal/cache"
	"github.com/Gary-Yez/go-admin/internal/state"
)

type localLink struct {
	Token     string
	ExpiresAt time.Time
}

func localLinkKey(id uint) string {
	return fmt.Sprintf("sys_file:link:%d", id)
}

func readLocalLink(key string) (localLink, error) {
	var link localLink
	if err := state.Cache().GetJSON(key, &link); err != nil {
		return link, err
	}
	if len(link.Token) != 64 || !time.Now().Before(link.ExpiresAt) {
		return localLink{}, cache.ErrCacheNotFound
	}
	return link, nil
}

// localFileLink 每个文件只缓存一个凭证，预览与下载共用，命中时不续期。
func localFileLink(ctx context.Context, id uint, preview bool) (string, int, error) {
	key := localLinkKey(id)
	link, err := readLocalLink(key)
	if errors.Is(err, cache.ErrCacheNotFound) {
		ctx, cancel := context.WithTimeout(ctx, 6*time.Second)
		defer cancel()
		lock, lockErr := state.Cache().WaitForLock(ctx, key, 15*time.Second, 5*time.Second)
		if lockErr != nil {
			return "", 0, lockErr
		}
		defer lock.Unlock()
		// 并发请求等待后重新读取，避免为同一文件生成多个凭证。
		link, err = readLocalLink(key)
		if errors.Is(err, cache.ErrCacheNotFound) {
			ttl, ttlErr := readLinkTTL()
			if ttlErr != nil {
				return "", 0, ttlErr
			}
			token, tokenErr := randomKey()
			if tokenErr != nil {
				return "", 0, tokenErr
			}
			link = localLink{Token: token, ExpiresAt: time.Now().Add(ttl)}
			err = state.Cache().SetJSON(key, link, ttl)
		}
	}
	if err != nil {
		return "", 0, err
	}
	remaining := time.Until(link.ExpiresAt)
	if remaining <= 0 {
		return "", 0, errors.New("链接已过期，请重新获取")
	}
	seconds := int((remaining + time.Second - 1) / time.Second)
	url := fmt.Sprintf("/sys_file/content/%d.%s", id, link.Token)
	if !preview {
		url += "?download=1"
	}
	return url, seconds, nil
}

// 地址中的文件编号仅用于定位缓存，实际访问必须核对随机凭证。
func resolveLocalLink(value string) (uint, error) {
	invalid := errors.New("下载链接无效或已过期")
	parts := strings.Split(value, ".")
	if len(parts) != 2 || len(parts[1]) != 64 {
		return 0, invalid
	}
	id, err := strconv.ParseUint(parts[0], 10, strconv.IntSize)
	if err != nil || id == 0 {
		return 0, invalid
	}
	link, err := readLocalLink(localLinkKey(uint(id)))
	if err != nil {
		return 0, err
	}
	if subtle.ConstantTimeCompare([]byte(parts[1]), []byte(link.Token)) != 1 {
		return 0, invalid
	}
	return uint(id), nil
}
