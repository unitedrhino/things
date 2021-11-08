package svc

import (
	"context"
	"gitee.com/godLei6/things/shared/db/mongodb"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/dmsvr/device"
	"gitee.com/godLei6/things/src/dmsvr/internal/config"
	"gitee.com/godLei6/things/src/dmsvr/internal/repo/model"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

type ServiceContext struct {
	Config      config.Config
	DeviceInfo  model.DeviceInfoModel
	ProductInfo model.ProductInfoModel
	DeviceLog   model.DeviceLogModel
	DmDB        model.DmModel
	DeviceID    *utils.SnowFlake
	ProductID   *utils.SnowFlake
	DevClient   *device.DevClient
	Mongo       *mongo.Database
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	di := model.NewDeviceInfoModel(conn, c.CacheRedis)
	pi := model.NewProductInfoModel(conn, c.CacheRedis)
	dl := model.NewDeviceLogModel(conn)
	DmDB := model.NewDmModel(conn, c.CacheRedis)
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
	return &ServiceContext{
		Config:      c,
		DeviceInfo:  di,
		ProductInfo: pi,
		DmDB:        DmDB,
		DeviceID:    DeviceID,
		ProductID:   ProductID,
		DeviceLog:   dl,
		DevClient:   devClient,
		Mongo:       mongoDB,
	}
}
