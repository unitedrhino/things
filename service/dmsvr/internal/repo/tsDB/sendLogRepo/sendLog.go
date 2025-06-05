package sendLogRepo

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceGroup"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/tdengine/sendLogRepo"
	"github.com/zeromicro/go-zero/core/logx"
)

type Send struct {
	*deviceLog.Send
}

func (m *Send) TableName() string {
	return "dm_time_send_log"
}

type SendLogRepo struct {
	db          *stores.DB
	asyncInsert *stores.AsyncInsert[Send]
}

func NewSendLogRepo(dataSource conf.TSDB, g []*deviceGroup.GroupDetail) deviceLog.SendRepo {
	if dataSource.DBType == conf.Tdengine {
		return sendLogRepo.NewSendLogRepo(dataSource, g)
	}
	stores.InitTsConn(dataSource)
	db := stores.GetTsConn(context.Background())
	var NeedInitColumn bool
	if db.Migrator().HasTable(&Send{}) {
		//需要初始化表
		NeedInitColumn = true
	}
	err := db.AutoMigrate(&Send{})
	logx.Must(err)
	if NeedInitColumn && stores.GetTsDBType() == conf.Pgsql {
		db.Exec("SELECT create_hypertable('dm_time_send_log','ts', chunk_time_interval => interval '1 day'    );")
	}
	return &SendLogRepo{db: db, asyncInsert: stores.NewAsyncInsert[Send](db, "")}
}
