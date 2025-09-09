package dmExport

import (
	"context"
	"encoding/json"

	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/things/service/dmsvr/client/devicemanage"
	"gitee.com/unitedrhino/things/service/dmsvr/client/productmanage"
	"gitee.com/unitedrhino/things/service/dmsvr/client/userdevice"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/userShared"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"gitee.com/unitedrhino/things/share/topics"
)

type ProductCacheT = *caches.Cache[dm.ProductInfo, string]

func NewProductInfoCache(pm productmanage.ProductManage, fastEvent *eventBus.FastEvent) (ProductCacheT, error) {
	return caches.NewCache(caches.CacheConfig[dm.ProductInfo, string]{
		KeyType:   topics.ServerCacheKeyDmProduct,
		FastEvent: fastEvent,
		GetData: func(ctx context.Context, key string) (*dm.ProductInfo, error) {
			ret, err := pm.ProductInfoRead(ctx, &dm.ProductInfoReadReq{ProductID: key, WithCategory: true, WithProtocol: true})
			return ret, err
		},
	})
}

type DeviceCacheT = *caches.Cache[dm.DeviceInfo, devices.Core]

func NewDeviceInfoCache(devM devicemanage.DeviceManage, fastEvent *eventBus.FastEvent) (DeviceCacheT, error) {
	return caches.NewCache(caches.CacheConfig[dm.DeviceInfo, devices.Core]{
		KeyType:   topics.ServerCacheKeyDmDevice,
		FastEvent: fastEvent,
		GetData: func(ctx context.Context, key devices.Core) (*dm.DeviceInfo, error) {
			ret, err := devM.DeviceInfoRead(ctx, &dm.DeviceInfoReadReq{ProductID: key.ProductID, DeviceName: key.DeviceName})
			return ret, err
		},
	})
}

//	type UserShareKey struct {
//		ProductID  string `json:"productID"`  //产品id
//		DeviceName string `json:"deviceName"` //设备名称
//		SharedUserID int64 `json:"sharedUserID"`
//	}
type UserShareCacheT = *caches.Cache[dm.UserDeviceShareInfo, userShared.UserShareKey]

func NewUserShareCache(devM userdevice.UserDevice, fastEvent *eventBus.FastEvent) (UserShareCacheT, error) {
	return caches.NewCache(caches.CacheConfig[dm.UserDeviceShareInfo, userShared.UserShareKey]{
		KeyType:   topics.ServerCacheKeyDmUserShareDevice,
		FastEvent: fastEvent,
		GetData: func(ctx context.Context, key userShared.UserShareKey) (*dm.UserDeviceShareInfo, error) {
			ret, err := devM.UserDeviceShareRead(ctx, &dm.UserDeviceShareReadReq{
				Device: &dm.DeviceCore{
					ProductID:  key.ProductID,
					DeviceName: key.DeviceName,
				},
			})
			return ret, err
		},
	})
}

type ProductSchemaCacheT = *caches.Cache[schema.Model, string]

func NewProductSchemaCache(pm productmanage.ProductManage, fastEvent *eventBus.FastEvent) (ProductSchemaCacheT, error) {
	return caches.NewCache(caches.CacheConfig[schema.Model, string]{
		KeyType:   topics.ServerCacheKeyDmProductSchema,
		FastEvent: fastEvent,
		Fmt: func(ctx context.Context, key string, data *schema.Model) *schema.Model {
			data.ValidateWithFmt()
			return data
		},
		GetData: func(ctx context.Context, key string) (*schema.Model, error) {
			if key == "" {
				return nil, errors.Parameter.AddMsgf("产品ID必填")
			}
			info, err := pm.ProductSchemaTslRead(ctx, &dm.ProductSchemaTslReadReq{ProductID: key})
			if err != nil {
				return nil, err
			}
			return schema.ValidateWithFmt([]byte(info.Tsl))
		},
	})
}

type DeviceSchemaCacheT = *caches.Cache[schema.Model, devices.Core]

func NewDeviceSchemaCache(pm devicemanage.DeviceManage, pc ProductSchemaCacheT, fastEvent *eventBus.FastEvent) (DeviceSchemaCacheT, error) {
	ret, err := caches.NewCache(caches.CacheConfig[schema.Model, devices.Core]{
		KeyType:   topics.ServerCacheKeyDmDeviceSchema,
		FastEvent: fastEvent,
		Fmt: func(ctx context.Context, key devices.Core, data *schema.Model) *schema.Model {
			pd, _ := pc.GetData(ctx, key.ProductID)
			newOne := data.Copy().Aggregation(pd)
			newOne.ValidateWithFmt()
			return newOne
		},
		GetData: func(ctx context.Context, key devices.Core) (*schema.Model, error) {
			if key.ProductID == "" || key.DeviceName == "" {
				return nil, errors.Parameter.AddMsgf("产品ID和设备ID必填")
			}
			info, err := pm.DeviceSchemaTslRead(ctx, &dm.DeviceSchemaTslReadReq{ProductID: key.ProductID, DeviceName: key.DeviceName})
			if err != nil {
				return nil, err
			}
			return schema.ValidateWithFmt([]byte(info.Tsl))
		},
	})
	if err != nil {
		return nil, err
	}
	pc.AddNotifySlot(func(ctx context.Context, keyB []byte) {
		var pKey devices.Core
		json.Unmarshal(keyB, &pKey)
		productID := string(keyB)
		ret.DeleteByFunc(func(key string) bool {
			ck := devices.Core{}
			json.Unmarshal([]byte(key), &ck)
			if ck.ProductID == productID {
				return true
			}
			return false
		})
	})
	return ret, nil
}
