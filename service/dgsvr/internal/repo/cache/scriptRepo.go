package cache

import (
	"context"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/things/service/dgsvr/internal/domain/custom"
	"github.com/dgraph-io/ristretto"
	"time"
)

const (
	expireTime = time.Hour
)

type (
	GetScriptInfo func(ctx context.Context, productID string) (info *custom.Info, err error)
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
func (t ScriptRepo) GetProtoFunc(ctx context.Context, productID string, dir custom.ConvertType,
	handle string, /*对应 mqtt topic的第一个 thing ota config 等等*/
	Type string /*操作类型 从topic中提取 物模型下就是   property属性 event事件 action行为*/) (custom.ConvertFunc, error) {
	vm, err := t.getScriptVm(ctx, productID)
	if err != nil || vm == nil {
		return nil, err
	}
	if dir == custom.ConvertTypeDown {
		fun := vm.DataDown(ctx, handle, Type)
		return fun, nil
	}
	fun := vm.DataUp(ctx, handle, Type)
	return fun, nil
}

func (t ScriptRepo) GetTransFormFunc(ctx context.Context, productID string, direction devices.Direction) (custom.TransFormFunc, error) {
	vm, err := t.getScriptVm(ctx, productID)
	if err != nil || vm == nil {
		return nil, err
	}
	return vm.TransformPayload(ctx, direction), nil
}

func (t ScriptRepo) getScriptVm(ctx context.Context, productID string) (*custom.Vm, error) {
	temp, ok := t.cache.Get(productID)
	if ok {
		if temp == nil {
			return nil, nil
		}
		return temp.(*custom.Vm), nil
	}
	info, err := t.getScript(ctx, productID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		t.cache.SetWithTTL(productID, nil, 1, expireTime)
		return nil, nil
	}
	vm := info.InitScript()
	t.cache.SetWithTTL(productID, vm, 1, expireTime)
	return vm, nil
}

func (t ScriptRepo) ClearCache(ctx context.Context, productID string) error {
	t.cache.Del(productID)
	info, err := t.getScript(ctx, productID)
	if err != nil {
		return err
	}
	t.cache.SetWithTTL(productID, info.InitScript(), 1, expireTime)
	return nil
}
