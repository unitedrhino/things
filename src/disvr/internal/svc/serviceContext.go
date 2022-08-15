package svc

import (
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/src/disvr/internal/config"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/repo/event/publish/pubDev"
	"github.com/i-Things/things/src/disvr/internal/repo/mysql"
	"github.com/i-Things/things/src/disvr/internal/repo/tdengine/hubLogRepo"
	"github.com/i-Things/things/src/disvr/internal/repo/tdengine/schemaDataRepo"
	"github.com/i-Things/things/src/disvr/internal/repo/tdengine/sdkLogRepo"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"os"
)

type ServiceContext struct {
	Config        config.Config
	DeviceInfo    mysql.DeviceInfoModel
	ProductInfo   mysql.ProductInfoModel
	PubDev        pubDev.PubDev
	SchemaMsgRepo deviceMsg.SchemaDataRepo
	HubLogRepo    deviceMsg.HubLogRepo
	SchemaRepo    schema.SchemaRepo
	SDKLogRepo    deviceMsg.SDKLogRepo
}

func NewServiceContext(c config.Config) *ServiceContext {
	hubLog := hubLogRepo.NewHubLogRepo(c.TDengine.DataSource)
	sdkLog := sdkLogRepo.NewSDKLogRepo(c.TDengine.DataSource)

	//TestTD(td)
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	di := mysql.NewDeviceInfoModel(conn)
	pi := mysql.NewProductInfoModel(conn)
	pt := mysql.NewProductSchemaModel(conn)
	tr := mysql.NewSchemaRepo(pt)
	deviceData := schemaDataRepo.NewSchemaDataRepo(c.TDengine.DataSource, tr.GetSchemaModel)
	pd, err := pubDev.NewPubDev(c.Event)
	if err != nil {
		logx.Error("NewPubDev err", err)
		os.Exit(-1)
	}
	return &ServiceContext{
		Config:        c,
		DeviceInfo:    di,
		ProductInfo:   pi,
		SchemaRepo:    tr,
		PubDev:        pd,
		SchemaMsgRepo: deviceData,
		HubLogRepo:    hubLog,
		SDKLogRepo:    sdkLog,
	}
}
