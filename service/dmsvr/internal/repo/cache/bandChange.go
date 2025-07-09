package cache

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/things/share/devices"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/kv"
)

func genBindChangeKey(dev devices.Core) string {
	return fmt.Sprintf("dm:device:bindChange:%s:%s", dev.ProductID, dev.DeviceName)
}

type BindChange struct {
	cache kv.Store
}

func NewBindChange(store kv.Store) *BindChange {
	return &BindChange{cache: store}
}

func (d *BindChange) Set(ctx context.Context, dev devices.Core, projectID int64) error {
	err := d.cache.SetexCtx(ctx, genBindChangeKey(dev), cast.ToString(projectID), 240)
	return err
}

func (d *BindChange) Get(ctx context.Context, dev devices.Core) int64 {
	v, err := d.cache.GetCtx(ctx, genBindChangeKey(dev))
	if err != nil {
		return 0
	}
	return cast.ToInt64(v)
}
