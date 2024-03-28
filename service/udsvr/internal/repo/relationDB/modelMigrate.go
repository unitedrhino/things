package relationDB

import (
	"context"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/stores"
)

func Migrate(c conf.Database) error {
	if c.IsInitTable == false {
		return nil
	}
	db := stores.GetCommonConn(context.TODO())
	return db.AutoMigrate(
		&UdSceneInfo{},
		&UdSceneThenAction{},
		&UdDeviceTimerInfo{},
		&UdSceneIfTrigger{},
		&UdOpsWorkOrder{},
	)
}
