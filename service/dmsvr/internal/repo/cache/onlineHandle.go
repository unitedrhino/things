package cache

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceStatus"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"time"
)

func genOnlineCacheKey() string {
	return "dm:device:online:score"
}
func genOnlineLockKey() string {
	return "dm:device:online:lock"
}

type DeviceStatus struct {
	cache kv.Store
}

func NewDeviceStatus(store kv.Store) *DeviceStatus {
	return &DeviceStatus{cache: store}
}

func (d *DeviceStatus) Lock(ctx context.Context) (bool, error) {
	ok, err := d.cache.SetnxExCtx(ctx, genOnlineLockKey(), time.Now().Format("2006-01-02 15:04:05.999"), 30)
	return ok, err
}

func (d *DeviceStatus) UnLock(ctx context.Context) error {
	_, err := d.cache.DelCtx(ctx, genOnlineLockKey())
	return err
}

func (d *DeviceStatus) AddDevice(ctx context.Context, in *deviceStatus.ConnectMsg) error {
	_, err := d.cache.ZaddCtx(ctx, genOnlineCacheKey(), in.Timestamp.UnixMilli(), utils.MarshalNoErr(in))
	return err
}

func (d *DeviceStatus) DelDevices(ctx context.Context, devs ...*deviceStatus.ConnectMsg) error {
	vals, err := utils.MarshalSlices(devs)
	if err != nil {
		return err
	}
	_, err = d.cache.ZremCtx(ctx, genOnlineCacheKey(), utils.ToAnySlice(vals))
	return err
}

func (d *DeviceStatus) GetDevices(ctx context.Context) ([]*deviceStatus.ConnectMsg, error) {
	now := time.Now()
	//剔除在该时间内同时登录及登出的,然后将5s之前的数据入库
	vals, err := d.cache.ZrangeCtx(ctx, genOnlineCacheKey(), 0, now.Add(-time.Second*2).UnixMilli())
	if err != nil {
		return nil, err
	}
	return utils.UnmarshalSlices[deviceStatus.ConnectMsg](vals)
}
