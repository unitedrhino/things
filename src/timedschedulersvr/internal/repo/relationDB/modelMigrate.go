package relationDB

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm/clause"
)

func Migrate(c conf.Database) error {
	if c.IsInitTable == false {
		return nil
	}
	db := stores.GetCommonConn(nil)
	var needInitColumn bool
	if !db.Migrator().HasTable(&TimedTask{}) {
		//需要初始化表
		needInitColumn = true
	}
	err := db.AutoMigrate(
		&TimedTask{},
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
	MigrateTimedQueueJob = []TimedTask{
		{
			Group:          "order",
			Type:           "queue",
			SubType:        "natsJs",
			Name:           "订单检查",
			Code:           "order_check",
			Params:         `{"topic":"server.435","payload":"adfgawe"}`,
			CronExpression: "@every 2s",
			Status:         TaskStatusStop,
			Priority:       "critical",
		},
		{
			Group:          "order",
			Type:           "queue",
			SubType:        "natsJs",
			Name:           "订单检查2",
			Code:           "order_check2",
			Params:         `{"topic":"server.123","payload":"sdfgarg"}`,
			CronExpression: "@every 2s",
			Status:         TaskStatusStop,
			Priority:       "low",
		},
		{
			Group:          "sql",
			Type:           "sql",
			SubType:        "normal",
			Name:           "sql执行",
			Code:           "sql_Exec",
			Params:         `{"sql":"insert into test_table(name) values('123123') "}`,
			CronExpression: "@every 2s",
			Status:         TaskStatusPause,
			Priority:       "low",
		},
	}
)
