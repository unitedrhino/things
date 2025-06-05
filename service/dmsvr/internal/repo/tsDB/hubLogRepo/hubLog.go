package hubLogRepo

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceGroup"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/tdengine/hubLogRepo"
	"github.com/zeromicro/go-zero/core/logx"
)

type Hub struct {
	*deviceLog.Hub
}

func (m *Hub) TableName() string {
	return "dm_time_hub_log"
}

type HubLogRepo struct {
	db          *stores.DB
	asyncInsert *stores.AsyncInsert[Hub]
}

func NewHubLogRepo(dataSource conf.TSDB, g []*deviceGroup.GroupDetail) deviceLog.HubRepo {
	if dataSource.DBType == conf.Tdengine {
		return hubLogRepo.NewHubLogRepo(dataSource)
	}
	stores.InitTsConn(dataSource)
	db := stores.GetTsConn(context.Background())
	var NeedInitColumn bool
	if db.Migrator().HasTable(&Hub{}) {
		//需要初始化表
		NeedInitColumn = true
	}
	err := db.AutoMigrate(&Hub{})
	logx.Must(err)
	if NeedInitColumn && stores.GetTsDBType() == conf.Pgsql {
		db.Exec("SELECT create_hypertable('dm_time_hub_log','ts', chunk_time_interval => interval '1 day'    );")
	}
	return &HubLogRepo{db: db, asyncInsert: stores.NewAsyncInsert[Hub](db, "")}
}
