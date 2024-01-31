package cache

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/shared/devices"
	"gitee.com/i-Things/core/shared/errors"
	"github.com/dgraph-io/ristretto"
	"github.com/dop251/goja"
	"github.com/i-Things/things/src/dmsvr/internal/domain/productCustom"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"sync"
	"time"
)

const (
	expireTime = time.Hour
)

type (
	ScriptRepo struct {
		cache *ristretto.Cache
	}
)
type CustomCacheStu struct {
	LoginAuthPool *sync.Pool
	Topics        []*productCustom.CustomTopic
}

func NewScriptRepo() *ScriptRepo {
	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	return &ScriptRepo{
		cache: cache,
	}
}

func (t *ScriptRepo) GetCustomTopic(ctx context.Context, productID string) (topics []*productCustom.CustomTopic, err error) {
	c, err := t.getCache(ctx, productID)
	return c.Topics, err
}

func (t ScriptRepo) GetTransFormFunc(ctx context.Context, productID string) (productCustom.LoginAuthFunc, error) {
	c, err := t.getCache(ctx, productID)
	if err != nil || c == nil {
		return nil, err
	}
	vm := c.LoginAuthPool.Get().(*goja.Runtime)
	f, ok := goja.AssertFunction(vm.Get(productCustom.LoginAuthFuncName))
	if !ok {
		return nil, errors.Parameter.AddDetail("未找到函数:" + productCustom.LoginAuthFuncName)
	}
	return func(dir devices.Direction, clientID string, userName string, password string) (*devices.Core, error) {
		ret, err := f(goja.Undefined(), vm.ToValue(dir), vm.ToValue(clientID), vm.ToValue(userName), vm.ToValue(password))
		if err != nil {
			return nil, err
		}
		v, err := ret.ToObject(vm).MarshalJSON()
		if err != nil {
			return nil, err
		}
		fmt.Println(v)
		return nil, nil
	}, nil
}

func (t *ScriptRepo) getCache(ctx context.Context, productID string) (*CustomCacheStu, error) {
	temp, ok := t.cache.Get(productID)
	if ok {
		if temp == nil {
			return nil, nil
		}
		return temp.(*CustomCacheStu), nil
	}
	ps, err := relationDB.NewProductCustomRepo(ctx).FindOneByProductID(ctx, productID)
	if err != nil {
		if err == errors.NotFind {
			return nil, nil
		}
		return nil, err
	}
	ret := &CustomCacheStu{
		Topics: ps.CustomTopics,
	}
	if ps.LoginAuthScript != "" {
		ret.LoginAuthPool = &sync.Pool{New: func() any {
			vm := goja.New()
			_, err := vm.RunString(ps.LoginAuthScript)
			if err != nil {
				return nil
			}
			return vm
		},
		}
	}
	t.cache.SetWithTTL(productID, ps, 1, expireTime)
	return ret, nil
}

func (t *ScriptRepo) ClearCache(ctx context.Context, productID string) error {
	t.cache.Del(productID)
	info, err := t.getCache(ctx, productID)
	if err != nil {
		return err
	}
	t.cache.SetWithTTL(productID, info, 1, expireTime)
	return nil
}
