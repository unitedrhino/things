package svc

import (
	"context"
	"github.com/go-things/things/shared/db/mongodb"
	"github.com/go-things/things/shared/utils"
	"github.com/go-things/things/src/dmsvr/device"
	"github.com/go-things/things/src/dmsvr/internal/config"
	"github.com/go-things/things/src/dmsvr/internal/repo"
	"github.com/go-things/things/src/dmsvr/internal/repo/model/mongorepo"
	"github.com/go-things/things/src/dmsvr/internal/repo/model/mysql"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"os"
)

type ServiceContext struct {
	Config          config.Config
	DeviceInfo      mysql.DeviceInfoModel
	ProductInfo     mysql.ProductInfoModel
	ProductTemplate mysql.ProductTemplateModel
	DeviceLog       mysql.DeviceLogModel
	DmDB            mysql.DmModel
	DeviceID        *utils.SnowFlake
	ProductID       *utils.SnowFlake
	DevClient       *device.DevClient
	DeviceData      repo.GetDeviceDataRepo
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	di := mysql.NewDeviceInfoModel(conn, c.CacheRedis)
	pi := mysql.NewProductInfoModel(conn, c.CacheRedis)
	pt := mysql.NewProductTemplateModel(conn, c.CacheRedis)
	dl := mysql.NewDeviceLogModel(conn)
	DmDB := mysql.NewDmModel(conn, c.CacheRedis)
	DeviceID := utils.NewSnowFlake(c.NodeID)
	ProductID := utils.NewSnowFlake(c.NodeID)

	devClient := device.NewDevClient(c.DevClient)
	//if token := mc.Connect(); token.Wait() && token.Error() != nil {
	//	panic(fmt.Sprintf("mqtt client connect err:%s",token.Error()))
	//}
	//token := mc.Publish("21CYs1k9YpG/test8/54598", 0, false, clientID+" send msg")
	//token.Wait()
	//time.Sleep(time.Hour)
	mongoDB, err := mongodb.NewMongo(c.Mongo.Url, c.Mongo.Database, context.TODO())
	if err != nil {
		logx.Error(err)
		os.Exit(-1)
	}
	dd := mongorepo.NewDeviceDataRepo(mongoDB)
	return &ServiceContext{
		Config:          c,
		DeviceInfo:      di,
		ProductInfo:     pi,
		ProductTemplate: pt,
		DmDB:            DmDB,
		DeviceID:        DeviceID,
		ProductID:       ProductID,
		DeviceLog:       dl,
		DevClient:       devClient,
		DeviceData:      dd,
	}
}
