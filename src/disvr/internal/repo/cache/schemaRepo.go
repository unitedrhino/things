package cache

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"github.com/i-Things/things/shared/domain/schema"
	"time"
)

const (
	expireTime = time.Hour
)

type (
	GetSchemaInfo func(ctx context.Context, productID string) (info *schema.Model, err error)
	SchemaRepo    struct {
		getSchema GetSchemaInfo
		cache     *ristretto.Cache
	}
)

func NewSchemaRepo(t GetSchemaInfo) schema.ReadRepo {
	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	return &SchemaRepo{
		getSchema: t,
		cache:     cache,
	}
}

func (t SchemaRepo) GetSchemaInfo(ctx context.Context, productID string) (*schema.Model, error) {
	temp, err := t.getSchema(ctx, productID)
	if err != nil {
		return nil, err
	}
	return temp, nil
}

func (t SchemaRepo) GetSchemaModel(ctx context.Context, productID string) (*schema.Model, error) {
	temp, ok := t.cache.Get(productID)
	if ok {
		return temp.(*schema.Model), nil
	}
	schemaInfo, err := t.getSchema(ctx, productID)
	if err != nil {
		return nil, err
	}
	t.cache.SetWithTTL(productID, schemaInfo, 1, expireTime)
	return schemaInfo, nil
}

func (t SchemaRepo) ClearCache(ctx context.Context, productID string) error {
	t.cache.Del(productID)
	return nil
}
