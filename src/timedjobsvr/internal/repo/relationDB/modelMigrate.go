package relationDB

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/stores"
)

func Migrate(c conf.Database) error {
	if c.IsInitTable == false {
		return nil
	}
	db := stores.GetCommonConn(nil)
	var needInitColumn bool
	if !db.Migrator().HasTable(&TimedJobLog{}) {
		//需要初始化表
		needInitColumn = true
	}
	err := db.AutoMigrate(
		&TimedJobLog{},
	)
	if needInitColumn {
		return migrateTableColumn()
	}
	return err
}
func migrateTableColumn() error {
	//db := stores.GetCommonConn(nil).Clauses(clause.OnConflict{DoNothing: true})
	return nil
}
