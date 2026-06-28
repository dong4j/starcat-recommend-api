package provider

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// CachedProvider 为上游 recommend 结果加进程内 TTL cache。
//
// 首版服务没有持久化层, 进程内缓存足够降低 SimRepo 重复调用压力。未来如果需要多实例
// 共享缓存, 在这里替换为 Redis/SQLite, handler 和客户端契约不需要变化。
type CachedProvider struct {
	base       Provider
	successTTL time.Duration
	emptyTTL   time.Duration
	errorTTL   time.Duration

	mu    sync.RWMutex
	items map[string]cacheEntry
}

type cacheEntry struct {
	result    Result
	err       error
	expiresAt time.Time
}

func NewCachedProvider(base Provider, successTTL, emptyTTL, errorTTL time.Duration) *CachedProvider {
	return &CachedProvider{
		base:       base,
		successTTL: successTTL,
		emptyTTL:   emptyTTL,
		errorTTL:   errorTTL,
		items:      map[string]cacheEntry{},
	}
}

func (p *CachedProvider) Recommend(ctx context.Context, query Query) (Result, error) {
	key := cacheKey(query)
	now := time.Now()
	if entry, ok := p.get(key, now); ok {
		if entry.err != nil {
			return Result{}, entry.err
		}
		entry.result.CacheStatus = "hit"
		return entry.result, nil
	}

	result, err := p.base.Recommend(ctx, query)
	ttl := p.successTTL
	if err != nil {
		ttl = p.errorTTL
	} else if len(result.Response.Items) == 0 {
		ttl = p.emptyTTL
	}
	p.set(key, cacheEntry{
		result:    result,
		err:       err,
		expiresAt: now.Add(ttl),
	})
	return result, err
}

func (p *CachedProvider) get(key string, now time.Time) (cacheEntry, bool) {
	p.mu.RLock()
	entry, ok := p.items[key]
	p.mu.RUnlock()
	if !ok || now.After(entry.expiresAt) {
		return cacheEntry{}, false
	}
	return entry, true
}

func (p *CachedProvider) set(key string, entry cacheEntry) {
	p.mu.Lock()
	p.items[key] = entry
	p.mu.Unlock()
}

func cacheKey(query Query) string {
	return fmt.Sprintf("%d:%d:%d", query.RepoID, query.Limit, query.Offset)
}
