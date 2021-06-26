package svc

import (
	"fmt"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/dmsvr/internal/config"
	"gitee.com/godLei6/things/src/dmsvr/internal/exchange/types"
	"gitee.com/godLei6/things/src/dmsvr/model"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type ServiceContext struct {
	Config      config.Config
	DeviceInfo  model.DeviceInfoModel
	ProductInfo model.ProductInfoModel
	DeviceLog  model.DeviceLogModel
	DeviceID    *utils.SnowFlake
	ProductID   *utils.SnowFlake
	Mqtt 		mqtt.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	di := model.NewDeviceInfoModel(conn, c.CacheRedis)
	pi := model.NewProductInfoModel(conn, c.CacheRedis)
	dl := model.NewDeviceLogModel(conn)
	DeviceID := utils.NewSnowFlake(c.NodeID)
	ProductID := utils.NewSnowFlake(c.NodeID)

	opts := mqtt.NewClientOptions()
	for _,broker := range c.Mqtt.Brokers{
		opts.AddBroker(broker)
	}
	clientID := fmt.Sprintf("%s:%d",c.Name,c.NodeID)
	opts.SetClientID(clientID).SetUsername(c.Mqtt.User).
		SetPassword(c.Mqtt.Pass).SetAutoReconnect(true).SetConnectRetry(true)
	opts.OnConnect = func(client mqtt.Client) {
		logx.Info("Connected")
	}
	mc:= mqtt.NewClient(opts)
	mc.Connect()
	//if token := mc.Connect(); token.Wait() && token.Error() != nil {
	//	panic(fmt.Sprintf("mqtt client connect err:%s",token.Error()))
	//}
	//token := mc.Publish("21CYs1k9YpG/test8/54598", 0, false, clientID+" send msg")
	//token.Wait()
	//time.Sleep(time.Hour)
	return &ServiceContext{
		Config:      c,
		DeviceInfo:  di,
		ProductInfo: pi,
		DeviceID:    DeviceID,
		ProductID:   ProductID,
		DeviceLog: dl,
		Mqtt: mc,
	}
}

func (l *ServiceContext) LogHandle(msg *types.Elements) error {
	ld, err := dm.GetClientIDInfo(msg.ClientID)
	if err != nil {
		return err
	}
	_, err = l.DeviceLog.Insert(model.DeviceLog{
		ProductID:   ld.ProductID,
		Action:      msg.Action,
		Timestamp:   time.Unix(msg.Timestamp, 0), // 操作时间
		DeviceName:  ld.DeviceName,
		Payload:     msg.Payload,
		Topic:       msg.Topic,
		CreatedTime: time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}