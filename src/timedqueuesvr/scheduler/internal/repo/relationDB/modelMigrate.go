package relationDB

import (
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm/clause"
)

func Migrate() error {
	db := stores.GetCommonConn(nil)
	var needInitColumn bool
	if !db.Migrator().HasTable(&TimedQueueJob{}) {
		//需要初始化表
		needInitColumn = true
	}
	err := db.AutoMigrate(
		&TimedQueueJob{},
	)
	if needInitColumn {
		return migrateTableColumn()
	}
	return err
}
func migrateTableColumn() error {
	db := stores.GetCommonConn(nil).Clauses(clause.OnConflict{DoNothing: true})
	if err := db.CreateInBatches(&MigrateTimedQueueJob, 100).Error; err != nil {
		return err
	}
	return nil
}

var (
	MigrateTimedQueueJob = []TimedQueueJob{
		{
			Group:          "order",
			Type:           "queue",
			SubType:        "natsJs",
			Name:           "订单检查",
			Code:           "order_check",
			Params:         `{"topic":"timed.435","payload":"adfgawe"}`,
			CronExpression: "@every 2s",
			Status:         JobStatusPause,
			Priority:       "critical",
		},
		{
			Group:          "order",
			Type:           "natsJs",
			Name:           "订单检查2",
			Code:           "order_check2",
			Params:         `{"topic":"timed.123","payload":"sdfgarg"}`,
			CronExpression: "@every 2s",
			Status:         JobStatusPause,
			Priority:       "low",
		},
	}
)
