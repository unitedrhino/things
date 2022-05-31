package svc

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/config"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	"github.com/i-Things/things/src/dmsvr/internal/domain/service/deviceData"
	"github.com/i-Things/things/src/dmsvr/internal/domain/thing"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/dataUpdate"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/innerLink"
	mysql "github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/repo/tdengine/deviceDataRepo"
	"github.com/i-Things/things/src/dmsvr/internal/repo/tdengine/deviceDebugLogRepo"
	"github.com/i-Things/things/src/dmsvr/internal/repo/tdengine/deviceLogRepo"
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
	InnerLink      innerLink.InnerLink
	DataUpdate     dataUpdate.DataUpdate
	Store          kv.Store
	DeviceDataRepo deviceData.DeviceDataRepo
	DeviceLogRepo  device.LogRepo
	TemplateRepo   thing.TemplateRepo
	DebugLogRepo   device.DebugLogRepo
}

func NewServiceContext(c config.Config) *ServiceContext {
	deviceData := deviceDataRepo.NewDeviceDataRepo(c.TDengine.DataSource)
	deviceLog := deviceLogRepo.NewDeviceLogRepo(c.TDengine.DataSource)
	deviceDebugLog := deviceDebugLogRepo.NewDeviceDebugLogRepo(c.TDengine.DataSource)

	//TestTD(td)
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	di := mysql.NewDeviceInfoModel(conn, c.CacheRedis)
	pi := mysql.NewProductInfoModel(conn, c.CacheRedis)
	pt := mysql.NewProductTemplateModel(conn, c.CacheRedis)
	tr := mysql.NewTemplateRepo(pt)

	DmDB := mysql.NewDmModel(conn, c.CacheRedis)
	store := kv.NewStore(c.CacheRedis)
	nodeId := utils.GetNodeID(c.CacheRedis, c.Name)
	DeviceID := utils.NewSnowFlake(nodeId)
	ProductID := utils.NewSnowFlake(nodeId)
	il, err := innerLink.NewInnerLink(c.InnerLink)
	if err != nil {
		logx.Error("NewInnerLink err", err)
		os.Exit(-1)
	}
	du, err := dataUpdate.NewDataUpdate(c.InnerLink)
	if err != nil {
		logx.Error("NewDataUpdate err", err)
		os.Exit(-1)
	}
	return &ServiceContext{
		Config:         c,
		DeviceInfo:     di,
		ProductInfo:    pi,
		TemplateRepo:   tr,
		DmDB:           DmDB,
		DeviceID:       DeviceID,
		ProductID:      ProductID,
		InnerLink:      il,
		DataUpdate:     du,
		Store:          store,
		DeviceDataRepo: deviceData,
		DeviceLogRepo:  deviceLog,
		DebugLogRepo:   deviceDebugLog,
	}
}
