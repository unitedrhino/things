package abnormalLogRepo

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceGroup"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/tdengine/abnormalLogRepo"
	"github.com/zeromicro/go-zero/core/logx"
)

type Abnormal struct {
	*deviceLog.Abnormal
}

func (m *Abnormal) TableName() string {
	return "dm_time_abnormal_log"
}

type AbnormalLogRepo struct {
	db          *stores.DB
	asyncInsert *stores.AsyncInsert[Abnormal]
}

func NewAbnormalLogRepo(dataSource conf.TSDB, g []*deviceGroup.GroupDetail) deviceLog.AbnormalRepo {
	if dataSource.DBType == conf.Tdengine {
		return abnormalLogRepo.NewAbnormalLogRepo(dataSource, g)
	}
	stores.InitTsConn(dataSource)
	db := stores.GetTsConn(context.Background())
	var NeedInitColumn bool
	if db.Migrator().HasTable(&Abnormal{}) {
		//需要初始化表
		NeedInitColumn = true
	}
	err := db.AutoMigrate(&Abnormal{})
	logx.Must(err)
	if NeedInitColumn && stores.GetTsDBType() == conf.Pgsql {
		db.Exec("SELECT create_hypertable('dm_time_abnormal_log','ts', chunk_time_interval => interval '1 day'    );")
	}
	return &AbnormalLogRepo{db: db, asyncInsert: stores.NewAsyncInsert[Abnormal](db, "")}
}
