package svc

import (
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/eventBus"
	"github.com/i-Things/things/shared/oss"
	"github.com/i-Things/things/shared/store"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/config"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsgManage"
	"github.com/i-Things/things/src/dmsvr/internal/repo/cache"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/publish/dataUpdate"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/repo/tdengine/deviceDataRepo"
	"github.com/i-Things/things/src/dmsvr/internal/repo/tdengine/hubLogRepo"
	"github.com/i-Things/things/src/dmsvr/internal/repo/tdengine/sdkLogRepo"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"os"
)

type ServiceContext struct {
	Config         config.Config
	DeviceInfo     mysql.DmDeviceInfoModel
	DeviceID       *utils.SnowFlake
	ProductID      *utils.SnowFlake
	DataUpdate     dataUpdate.DataUpdate
	Cache          kv.Store
	SchemaManaRepo deviceMsgManage.SchemaDataRepo
	HubLogRepo     deviceMsgManage.HubLogRepo
	SchemaRepo     schema.Repo
	SDKLogRepo     deviceMsgManage.SDKLogRepo
	GroupID        *utils.SnowFlake
	Gateway        mysql.DmGatewayDeviceModel
	OssClient      *oss.Client
	Bus            eventBus.Bus
}

func NewServiceContext(c config.Config) *ServiceContext {
	hubLog := hubLogRepo.NewHubLogRepo(c.TDengine.DataSource)
	sdkLog := sdkLogRepo.NewSDKLogRepo(c.TDengine.DataSource)

	//TestTD(td)
	conn := sqlx.NewMysql(c.Database.DSN)
	di := mysql.NewDmDeviceInfoModel(conn)
	tr := cache.NewSchemaRepo()
	deviceData := deviceDataRepo.NewDeviceDataRepo(c.TDengine.DataSource, tr.GetSchemaModel)
	cache := kv.NewStore(c.CacheRedis)
	nodeId := utils.GetNodeID(c.CacheRedis, c.Name)
	DeviceID := utils.NewSnowFlake(nodeId)
	ProductID := utils.NewSnowFlake(nodeId)
	du, err := dataUpdate.NewDataUpdate(c.Event)

	if err != nil {
		logx.Error("NewDataUpdate err", err)
		os.Exit(-1)
	}
	GroupID := utils.NewSnowFlake(nodeId)
	gw := mysql.NewDmGatewayDeviceModel(conn)
	ossClient := oss.NewOssClient(c.OssConf)
	if ossClient == nil {
		logx.Error("NewOss err")
		os.Exit(-1)
	}
	bus := eventBus.NewEventBus()
	store.InitConn(c.Database)
	return &ServiceContext{
		Bus:            bus,
		Config:         c,
		OssClient:      ossClient,
		DeviceInfo:     di,
		SchemaRepo:     tr,
		DeviceID:       DeviceID,
		ProductID:      ProductID,
		DataUpdate:     du,
		Cache:          cache,
		SchemaManaRepo: deviceData,
		HubLogRepo:     hubLog,
		SDKLogRepo:     sdkLog,
		GroupID:        GroupID,
		Gateway:        gw,
	}
}
