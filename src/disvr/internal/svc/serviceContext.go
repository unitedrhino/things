package svc

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/src/disvr/internal/config"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/repo/cache"
	"github.com/i-Things/things/src/disvr/internal/repo/event/publish/pubDev"
	"github.com/i-Things/things/src/disvr/internal/repo/tdengine/hubLogRepo"
	"github.com/i-Things/things/src/disvr/internal/repo/tdengine/schemaDataRepo"
	"github.com/i-Things/things/src/disvr/internal/repo/tdengine/sdkLogRepo"
	devicemanage "github.com/i-Things/things/src/dmsvr/client/devicemanage"
	productmanage "github.com/i-Things/things/src/dmsvr/client/productmanage"
	"github.com/i-Things/things/src/dmsvr/dmdirect"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
	"time"
)

type ServiceContext struct {
	Config        config.Config
	PubDev        pubDev.PubDev
	SchemaMsgRepo deviceMsg.SchemaDataRepo
	HubLogRepo    deviceMsg.HubLogRepo
	SchemaRepo    schema.ReadRepo
	SDKLogRepo    deviceMsg.SDKLogRepo
	DeviceM       devicemanage.DeviceManage
	ProductM      productmanage.ProductManage
}

func NewServiceContext(c config.Config) *ServiceContext {
	var (
		deviceM  devicemanage.DeviceManage
		productM productmanage.ProductManage
	)

	hubLog := hubLogRepo.NewHubLogRepo(c.TDengine.DataSource)
	sdkLog := sdkLogRepo.NewSDKLogRepo(c.TDengine.DataSource)

	//TestTD(td)
	pd, err := pubDev.NewPubDev(c.Event)
	if err != nil {
		logx.Error("NewPubDev err", err)
		os.Exit(-1)
	}
	if c.DmRpc.Mode == conf.ClientModeGrpc {
		deviceM = devicemanage.NewDeviceManage(zrpc.MustNewClient(c.DmRpc.Conf))
		productM = productmanage.NewProductManage(zrpc.MustNewClient(c.DmRpc.Conf))
	} else {
		deviceM = dmdirect.NewDeviceManage(nil)
		productM = dmdirect.NewProductManage(nil)
	}
	tr := cache.NewSchemaRepo(func(ctx context.Context, productID string) (schema.Info, error) {
		info, err := productM.ProductSchemaRead(ctx, &dm.ProductSchemaReadReq{ProductID: productID})
		if err != nil {
			return schema.Info{}, err
		}
		return schema.Info{
			Schema:      info.Schema,
			ProductID:   productID,
			CreatedTime: time.Unix(info.CreatedTime, 0),
		}, nil
	})
	deviceData := schemaDataRepo.NewSchemaDataRepo(c.TDengine.DataSource, tr.GetSchemaModel)

	return &ServiceContext{
		Config:        c,
		SchemaRepo:    tr,
		PubDev:        pd,
		SchemaMsgRepo: deviceData,
		HubLogRepo:    hubLog,
		SDKLogRepo:    sdkLog,
		ProductM:      productM,
		DeviceM:       deviceM,
	}
}
