package svc

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/disvr/internal/config"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgHubLog"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgSdkLog"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgThing"
	"github.com/i-Things/things/src/disvr/internal/repo/event/publish/pubApp"
	"github.com/i-Things/things/src/disvr/internal/repo/event/publish/pubDev"
	"github.com/i-Things/things/src/disvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/disvr/internal/repo/tdengine/hubLogRepo"
	"github.com/i-Things/things/src/disvr/internal/repo/tdengine/schemaDataRepo"
	"github.com/i-Things/things/src/disvr/internal/repo/tdengine/sdkLogRepo"
	devicemanage "github.com/i-Things/things/src/dmsvr/client/devicemanage"
	productmanage "github.com/i-Things/things/src/dmsvr/client/productmanage"
	remoteconfig "github.com/i-Things/things/src/dmsvr/client/remoteconfig"
	"github.com/i-Things/things/src/dmsvr/dmdirect"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
)

type ServiceContext struct {
	Config        config.Config
	PubDev        pubDev.PubDev
	PubApp        pubApp.PubApp
	SchemaMsgRepo msgThing.SchemaDataRepo
	HubLogRepo    msgHubLog.HubLogRepo
	SchemaRepo    schema.ReadRepo
	SDKLogRepo    msgSdkLog.SDKLogRepo
	DeviceM       devicemanage.DeviceManage
	ProductM      productmanage.ProductManage
	RemoteConfig  remoteconfig.RemoteConfig
	Cache         kv.Store
}

func NewServiceContext(c config.Config) *ServiceContext {
	var (
		deviceM      devicemanage.DeviceManage
		productM     productmanage.ProductManage
		remoteConfig remoteconfig.RemoteConfig
	)

	hubLog := hubLogRepo.NewHubLogRepo(c.TDengine.DataSource)
	sdkLog := sdkLogRepo.NewSDKLogRepo(c.TDengine.DataSource)
	stores.InitConn(c.Database)
	err := relationDB.Migrate()
	if err != nil {
		logx.Error("disvr 数据库初始化失败 err", err)
		os.Exit(-1)
	}
	pd, err := pubDev.NewPubDev(c.Event)
	if err != nil {
		logx.Error("NewPubDev err", err)
		os.Exit(-1)
	}
	pa, err := pubApp.NewPubApp(c.Event)
	if err != nil {
		logx.Error("NewPubApp err", err)
		os.Exit(-1)
	}
	if c.DmRpc.Mode == conf.ClientModeGrpc {
		deviceM = devicemanage.NewDeviceManage(zrpc.MustNewClient(c.DmRpc.Conf))
		productM = productmanage.NewProductManage(zrpc.MustNewClient(c.DmRpc.Conf))
		remoteConfig = remoteconfig.NewRemoteConfig(zrpc.MustNewClient(c.DmRpc.Conf))
	} else {
		deviceM = dmdirect.NewDeviceManage(c.DmRpc.RunProxy)
		productM = dmdirect.NewProductManage(c.DmRpc.RunProxy)
		remoteConfig = dmdirect.NewRemoteConfig(c.DmRpc.RunProxy)
	}
	tr := schema.NewReadRepo(func(ctx context.Context, productID string) (*schema.Model, error) {
		info, err := productM.ProductSchemaTslRead(ctx, &dm.ProductSchemaTslReadReq{ProductID: productID})
		if err != nil {
			return nil, err
		}
		return schema.ValidateWithFmt([]byte(info.Tsl))
	})
	store := kv.NewStore(c.CacheRedis)

	deviceData := schemaDataRepo.NewSchemaDataRepo(c.TDengine.DataSource, tr.GetSchemaModel, store)
	return &ServiceContext{
		PubApp:        pa,
		Config:        c,
		SchemaRepo:    tr,
		PubDev:        pd,
		Cache:         store,
		SchemaMsgRepo: deviceData,
		HubLogRepo:    hubLog,
		SDKLogRepo:    sdkLog,
		ProductM:      productM,
		DeviceM:       deviceM,
		RemoteConfig:  remoteConfig,
	}
}
