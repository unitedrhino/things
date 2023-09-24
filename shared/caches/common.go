package caches

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"sync"
)

var (
	once  sync.Once
	store kv.Store
)

func InitStore(c cache.ClusterConf) {
	once.Do(func() { store = kv.NewStore(c) })
}
