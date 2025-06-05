package statusLogRepo

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceGroup"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/tdengine/statusLogRepo"
	"github.com/zeromicro/go-zero/core/logx"
)

type Status struct {
	*deviceLog.Status
}

func (m *Status) TableName() string {
	return "dm_time_status_log"
}

type StatusLogRepo struct {
	db          *stores.DB
	asyncInsert *stores.AsyncInsert[Status]
}

func NewStatusLogRepo(dataSource conf.TSDB, g []*deviceGroup.GroupDetail) deviceLog.StatusRepo {
	if dataSource.DBType == conf.Tdengine {
		return statusLogRepo.NewStatusLogRepo(dataSource)
	}
	stores.InitTsConn(dataSource)
	db := stores.GetTsConn(context.Background())
	var NeedInitColumn bool
	if db.Migrator().HasTable(&Status{}) {
		//需要初始化表
		NeedInitColumn = true
	}
	err := db.AutoMigrate(&Status{})
	logx.Must(err)
	if NeedInitColumn && stores.GetTsDBType() == conf.Pgsql {
		db.Exec("SELECT create_hypertable('dm_time_status_log','ts', chunk_time_interval => interval '1 day'    );")
	}
	return &StatusLogRepo{db: db, asyncInsert: stores.NewAsyncInsert[Status](db, "")}
}
