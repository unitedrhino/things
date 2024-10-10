package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
)

func Migrate(c conf.Database) error {
	if c.IsInitTable == false {
		return nil
	}
	db := stores.GetCommonConn(context.TODO())
	return db.AutoMigrate(
		&UdSceneLog{},
		&UdSceneInfo{},
		&UdSceneThenAction{},
		&UdDeviceTimerInfo{},
		&UdSceneIfTrigger{},
		&UdAlarmScene{},
		&UdAlarmRecord{},
		&UdAlarmInfo{},
	)
}
