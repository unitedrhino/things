package msvc

import (
	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/dmsvr/internal/config"
	"gitee.com/godLei6/things/src/dmsvr/internal/msgquque/types"
	"gitee.com/godLei6/things/src/dmsvr/model"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"time"
)

type ServiceContext struct {
	Config     config.Config
	DeviceLog  model.DeviceLogModel
	DeviceInfo model.DeviceInfoModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	di := model.NewDeviceInfoModel(conn, c.CacheRedis)
	dl := model.NewDeviceLogModel(conn)
	return &ServiceContext{
		Config:     c,
		DeviceLog:  dl,
		DeviceInfo: di,
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
