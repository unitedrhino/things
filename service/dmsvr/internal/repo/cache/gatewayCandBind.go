package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"time"
)

func genGatewayCanBindCacheKey(gateway devices.Core) string {
	return fmt.Sprintf("dm:device:gateway:canBind:%s:%s", gateway.ProductID, gateway.DeviceName)
}

type GatewayCanBind struct {
	cache kv.Store
}
type GatewayCanBindStu struct {
	Gateway     devices.Core
	SubDevices  []*devices.Core
	UpdatedTime int64 //秒时间戳
}

func NewGatewayCanBind(store kv.Store) *GatewayCanBind {
	return &GatewayCanBind{cache: store}
}

func (d *GatewayCanBind) Update(ctx context.Context, in *GatewayCanBindStu) error {
	err := d.cache.SetexCtx(ctx, genGatewayCanBindCacheKey(in.Gateway), utils.MarshalNoErr(in), int(time.Hour*24/time.Second))
	return err
}

func (d *GatewayCanBind) GetDevices(ctx context.Context, gateway devices.Core) (*GatewayCanBindStu, error) {
	vals, err := d.cache.Get(genGatewayCanBindCacheKey(gateway))
	if err != nil {
		return nil, err
	}
	if vals == "" {
		return nil, errors.NotFind
	}
	var ret GatewayCanBindStu
	err = json.Unmarshal([]byte(vals), &ret)
	return &ret, err
}
