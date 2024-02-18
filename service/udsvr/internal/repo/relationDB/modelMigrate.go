package relationDB

import (
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/stores"
)

func Migrate(c conf.Database) error {
	if c.IsInitTable == false {
		return nil
	}
	db := stores.GetCommonConn(nil)
	return db.AutoMigrate(
		&UdSceneInfo{},
		&UdDeviceTimingInfo{},
		&UdOpsWorkOrder{},
	)
}
