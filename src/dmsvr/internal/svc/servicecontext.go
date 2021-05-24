package svc

import (
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/dmsvr/internal/config"
	"gitee.com/godLei6/things/src/dmsvr/model"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	DeviceInfo 		model.DeviceInfoModel
	ProductInfo     model.ProductInfoModel
	DeviceID		*utils.SnowFlake
	ProductID		*utils.SnowFlake
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	di := model.NewDeviceInfoModel(conn,c.CacheRedis)
	pi := model.NewProductInfoModel(conn,c.CacheRedis)
	DeviceID := utils.NewSnowFlake(c.NodeID)
	ProductID := utils.NewSnowFlake(c.NodeID)
	return &ServiceContext{
		Config: c,
		DeviceInfo: di,
		ProductInfo: pi,
		DeviceID:DeviceID,
		ProductID: ProductID,
	}
}
