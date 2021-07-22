package svc

import (
	"context"
	"fmt"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/dmsvr/internal/config"
	"gitee.com/godLei6/things/src/dmsvr/model"
	"gitee.com/godLei6/things/src/dmsvr/mongodb"
	mqtt "github.com/eclipse/paho.mqtt.golang"
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
	Mqtt        mqtt.Client
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

	opts := mqtt.NewClientOptions()
	for _, broker := range c.Mqtt.Brokers {
		opts.AddBroker(broker)
	}
	clientID := fmt.Sprintf("%s:%d", c.Name, c.NodeID)
	opts.SetClientID(clientID).SetUsername(c.Mqtt.User).
		SetPassword(c.Mqtt.Pass).SetAutoReconnect(true).SetConnectRetry(true)
	opts.OnConnect = func(client mqtt.Client) {
		logx.Info("Connected")
	}
	mc := mqtt.NewClient(opts)
	mc.Connect()
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
		Mqtt:        mc,
		Mongo:       mongoDB,
	}
}
