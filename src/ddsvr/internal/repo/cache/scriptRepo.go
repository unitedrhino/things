package cache

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/src/ddsvr/internal/domain/script"
	"time"
)

const (
	expireTime = time.Hour
)

type (
	GetScriptInfo func(ctx context.Context, productID string) (info *script.Info, err error)
	ScriptRepo    struct {
		getScript GetScriptInfo
		cache     *ristretto.Cache
	}
)

func NewScriptRepo(t GetScriptInfo) *ScriptRepo {
	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	return &ScriptRepo{
		getScript: t,
		cache:     cache,
	}
}

// 自定义协议转iTHings协议
func (t ScriptRepo) GetRawDataToProtocol(ctx context.Context, productID string,
	handle string, /*对应 mqtt topic的第一个 thing ota config 等等*/
	Type string /*操作类型 从topic中提取 物模型下就是   property属性 event事件 action行为*/) (*schema.Model, error) {
	temp, ok := t.cache.Get(productID)
	if ok {
		return temp.(*schema.Model), nil
	}
	schemaInfo, err := t.getScript(ctx, productID)
	if err != nil {
		return nil, err
	}
	t.cache.SetWithTTL(productID, schemaInfo, 1, expireTime)
	return schemaInfo, nil
}

func (t ScriptRepo) ClearCache(ctx context.Context, productID string) error {
	t.cache.Del(productID)
	info, err := t.getScript(ctx, productID)
	if err != nil {
		return err
	}
	return nil
}
