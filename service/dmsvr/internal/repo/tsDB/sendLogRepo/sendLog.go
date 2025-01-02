package sendLogRepo

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/tdengine/sendLogRepo"
	"github.com/zeromicro/go-zero/core/logx"
)

type Send struct {
	*deviceLog.Send
}

func (m *Send) TableName() string {
	return "dm_send_log"
}

type SendLogRepo struct {
	db          *stores.DB
	asyncInsert *stores.AsyncInsert[Send]
}

func NewSendLogRepo(dataSource conf.TSDB) deviceLog.SendRepo {
	if dataSource.DBType == conf.Tdengine {
		return sendLogRepo.NewSendLogRepo(dataSource)
	}
	stores.InitTsConn(dataSource)
	db := stores.GetTsConn(context.Background())
	err := db.AutoMigrate(&Send{})
	logx.Must(err)
	return &SendLogRepo{db: db, asyncInsert: stores.NewAsyncInsert[Send](db, "")}
}
