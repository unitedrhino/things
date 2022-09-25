package svc

import (
	"github.com/i-Things/things/shared/domain/schema"
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
	DeviceInfo     mysql.DeviceInfoModel
	ProductInfo    mysql.ProductInfoModel
	DmDB           mysql.DmModel
	DeviceID       *utils.SnowFlake
	ProductID      *utils.SnowFlake
	DataUpdate     dataUpdate.DataUpdate
	Store          kv.Store
	SchemaManaRepo deviceMsgManage.SchemaDataRepo
	HubLogRepo     deviceMsgManage.HubLogRepo
	SchemaRepo     schema.Repo
	SDKLogRepo     deviceMsgManage.SDKLogRepo
	FirmwareInfo   mysql.ProductFirmwareModel
	GroupInfo      mysql.GroupInfoModel
	GroupDevice    mysql.GroupDeviceModel
	GroupID        *utils.SnowFlake
	GroupDB        mysql.GroupModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	hubLog := hubLogRepo.NewHubLogRepo(c.TDengine.DataSource)
	sdkLog := sdkLogRepo.NewSDKLogRepo(c.TDengine.DataSource)

	//TestTD(td)
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	di := mysql.NewDeviceInfoModel(conn)
	pi := mysql.NewProductInfoModel(conn)
	pt := mysql.NewProductSchemaModel(conn)
	tr := cache.NewSchemaRepo(pt)
	deviceData := deviceDataRepo.NewDeviceDataRepo(c.TDengine.DataSource, tr.GetSchemaModel)
	fr := mysql.NewProductFirmwareModel(conn)
	DmDB := mysql.NewDmModel(conn)
	store := kv.NewStore(c.CacheRedis)
	nodeId := utils.GetNodeID(c.CacheRedis, c.Name)
	DeviceID := utils.NewSnowFlake(nodeId)
	ProductID := utils.NewSnowFlake(nodeId)
	du, err := dataUpdate.NewDataUpdate(c.Event)
	gi := mysql.NewGroupInfoModel(conn)
	gd := mysql.NewGroupDeviceModel(conn)
	GroupID := utils.NewSnowFlake(nodeId)
	GroupDB := mysql.NewGroupModel(conn, c.CacheRedis)
	if err != nil {
		logx.Error("NewDataUpdate err", err)
		os.Exit(-1)
	}
	return &ServiceContext{
		Config:         c,
		DeviceInfo:     di,
		ProductInfo:    pi,
		FirmwareInfo:   fr,
		SchemaRepo:     tr,
		DmDB:           DmDB,
		DeviceID:       DeviceID,
		ProductID:      ProductID,
		DataUpdate:     du,
		Store:          store,
		SchemaManaRepo: deviceData,
		HubLogRepo:     hubLog,
		SDKLogRepo:     sdkLog,
		GroupInfo:      gi,
		GroupDevice:    gd,
		GroupID:        GroupID,
		GroupDB:        GroupDB,
	}
}
