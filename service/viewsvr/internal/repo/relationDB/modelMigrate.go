package relationDB

import (
	"gitee.com/i-Things/core/shared/conf"
	"gitee.com/i-Things/core/shared/stores"
)

func Migrate(c conf.Database) error {
	if c.IsInitTable == false {
		return nil
	}
	db := stores.GetCommonConn(nil)
	return db.AutoMigrate(
		&ViewProjectInfo{},
		&ViewProjectDetail{},
	)
}
