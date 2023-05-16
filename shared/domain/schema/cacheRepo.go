package schema

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"sync"
	"time"
)

const (
	expireTime = time.Hour
)

type (
	GetReadInfo   func(ctx context.Context, productID string) (info *Model, err error)
	CacheReadRepo struct {
		getSchema GetReadInfo
		cache     *ristretto.Cache
	}
)

var (
	cacheReadRepo *CacheReadRepo
	ReadOnce      sync.Once
)

func NewReadRepo(t GetReadInfo) ReadRepo {
	ReadOnce.Do(func() { //单体模式避免反复缓存
		cache, _ := ristretto.NewCache(&ristretto.Config{
			NumCounters: 1e7,     // number of keys to track frequency of (10M).
			MaxCost:     1 << 30, // maximum cost of cache (1GB).
			BufferItems: 64,      // number of keys per Get buffer.
		})
		cacheReadRepo = &CacheReadRepo{
			getSchema: t,
			cache:     cache,
		}
	})
	return cacheReadRepo
}

func (t CacheReadRepo) GetSchemaInfo(ctx context.Context, productID string) (*Model, error) {
	temp, err := t.getSchema(ctx, productID)
	if err != nil {
		return nil, err
	}
	return temp, nil
}

func (t CacheReadRepo) GetSchemaModel(ctx context.Context, productID string) (*Model, error) {
	temp, ok := t.cache.Get(productID)
	if ok {
		return temp.(*Model), nil
	}
	schemaInfo, err := t.getSchema(ctx, productID)
	if err != nil {
		return nil, err
	}
	t.cache.SetWithTTL(productID, schemaInfo, 1, expireTime)
	return schemaInfo, nil
}

func (t CacheReadRepo) ClearCache(ctx context.Context, productID string) error {
	t.cache.Del(productID)
	return nil
}
