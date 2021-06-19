package msvc

import (
	"gitee.com/godLei6/things/src/dmsvr/device/model"
	"gitee.com/godLei6/things/src/dmsvr/device/msgquque/config"
	"gitee.com/godLei6/things/src/dmsvr/device/msgquque/types"
	"gitee.com/godLei6/things/src/dmsvr/dm"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"time"
)

type ServiceContext struct {
	Config config.Config
	DeviceLog model.DeviceLogModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	dl := model.NewDeviceLogModel(conn)
	return &ServiceContext{
		Config: c,
		DeviceLog:dl,
	}
}


func (l *ServiceContext) LogHandle(msg *types.Elements) error {
	ld,err := dm.GetClientIDInfo(msg.ClientID)
	if err!= nil {
		return err
	}
	_,err = l.DeviceLog.Insert(model.DeviceLog{
		ProductID   :ld.ProductID,
		Action      :msg.Action,
		Timestamp   :time.Unix(msg.Timestamp,0),  // 操作时间
		DeviceName  :ld.DeviceName,
		Payload     :msg.Payload,
		Topic       :msg.Topic,
		CreatedTime :time.Now(),
	})
	if err!= nil {
		return err
	}
	return nil
}