package relationDB

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/timed/internal/domain"
	"gorm.io/gorm/clause"
	"sync"
)

var once sync.Once

func Migrate(c conf.Database) (err error) {
	if c.IsInitTable == false {
		return
	}
	once.Do(func() {
		db := stores.GetCommonConn(nil)
		var needInitColumn bool
		if !db.Migrator().HasTable(&TimedTaskGroup{}) {
			//需要初始化表
			needInitColumn = true
		}
		err = db.AutoMigrate(
			&TimedTaskLog{},
			&TimedTaskGroup{},
			&TimedTask{},
		)
		if err != nil {
			return
		}
		if needInitColumn {
			err = migrateTableColumn()
		}
	})
	return
}
func migrateTableColumn() error {
	db := stores.GetCommonConn(nil).Clauses(clause.OnConflict{DoNothing: true})
	if err := db.CreateInBatches(&MigrateTimedTask, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateTimedTaskGroup, 100).Error; err != nil {
		return err
	}
	return nil
}

var (
	MigrateTimedTask = []TimedTask{
		{
			GroupCode: "queueTest",
			Type:      domain.TaskTypeTiming,
			Name:      "消息发送",
			Code:      "msgSendTest",
			Params:    `{"topic":"server.435","payload":"adfgawe"}`,
			CronExpr:  "@every 2s",
			Status:    def.StatusWaitRun,
			Priority:  2,
		},
		{
			GroupCode: "sqlJsTest",
			Type:      domain.TaskTypeTiming,
			Name:      "脚本执行",
			Code:      "sqlExec",
			Params:    `{"execContent": "function SqlJob(){Set('123','sdafawef');let a = Get('123');LogInfo('get value:',a);let code = GetEnv('code');LogInfo('get code env:',code);Exec(\"insert into test_table(name) values('123123')\");let values = Select('select * from test_table limit 10');LogInfo('select get value :',values);return {code:200,msg:'ok'};}"}`,
			CronExpr:  "@every 2s",
			Status:    def.StatusWaitRun,
			Priority:  4,
		},
		{
			GroupCode: "queueTest",
			Type:      domain.TaskTypeDelay,
			Name:      "延时测试",
			Code:      "delayTest",
			Params:    `{"topic":"server.333","payload":"garegawef"}`,
			CronExpr:  "",
			Status:    def.StatusRunning,
			Priority:  3,
		},
		{
			GroupCode: "iThingsQueueTiming",
			Type:      domain.TaskTypeTiming,
			Name:      "job服务.sql类型.脚本执行.Redis hash缓存清理",
			Code:      "timedJobRedisHashClean",
			Params:    `{"topic":"server.timedjob.cache.hash.clean","payload":""}`,
			CronExpr:  "1 1 * * ?",
			Status:    def.StatusWaitRun,
			Priority:  3,
		},
	}
	MigrateTimedTaskGroup = []TimedTaskGroup{
		{
			Code:     "queueTest",
			Name:     "消息队列测试",
			Type:     "queue",
			SubType:  "natsJs",
			Priority: 9,
		},
		{
			Code:     "iThingsQueueTiming",
			Name:     "iThings系统定时消息任务组",
			Type:     "queue",
			SubType:  "natsJs",
			Priority: 9,
		},
		{
			Code:     "sqlJsTest",
			Name:     "sqlJs模式测试",
			Type:     "sql",
			SubType:  "js",
			Priority: 7,
			Env:      map[string]string{"code": "66666"},
			Config:   `{"database":{"select":{"dsn":"root:password@tcp(127.0.0.1:3306)/iThings?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai","dbType":"mysql"}}}`,
		},
	}
)
