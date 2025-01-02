package statusLogRepo

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/tdengine/statusLogRepo"
	"github.com/zeromicro/go-zero/core/logx"
)

type Status struct {
	*deviceLog.Status
}

func (m *Status) TableName() string {
	return "dm_status_log"
}

type StatusLogRepo struct {
	db          *stores.DB
	asyncInsert *stores.AsyncInsert[Status]
}

func NewStatusLogRepo(dataSource conf.TSDB) deviceLog.StatusRepo {
	if dataSource.DBType == conf.Tdengine {
		return statusLogRepo.NewStatusLogRepo(dataSource)
	}
	stores.InitTsConn(dataSource)
	db := stores.GetTsConn(context.Background())
	err := db.AutoMigrate(&Status{})
	logx.Must(err)
	return &StatusLogRepo{db: db, asyncInsert: stores.NewAsyncInsert[Status](db, "")}
}
