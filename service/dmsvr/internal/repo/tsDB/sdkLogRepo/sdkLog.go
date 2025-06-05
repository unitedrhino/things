package sdkLogRepo

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceGroup"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/tdengine/sdkLogRepo"
	"github.com/zeromicro/go-zero/core/logx"
)

type SDK struct {
	*deviceLog.SDK
}

func (m *SDK) TableName() string {
	return "dm_time_sdk_log"
}

type SDKLogRepo struct {
	db          *stores.DB
	asyncInsert *stores.AsyncInsert[SDK]
}

func NewSDKLogRepo(dataSource conf.TSDB, g []*deviceGroup.GroupDetail) deviceLog.SDKRepo {
	if dataSource.DBType == conf.Tdengine {
		return sdkLogRepo.NewSDKLogRepo(dataSource)
	}
	stores.InitTsConn(dataSource)
	db := stores.GetTsConn(context.Background())
	var NeedInitColumn bool
	if db.Migrator().HasTable(&SDK{}) {
		//需要初始化表
		NeedInitColumn = true
	}
	err := db.AutoMigrate(&SDK{})
	logx.Must(err)
	if NeedInitColumn && stores.GetTsDBType() == conf.Pgsql {
		db.Exec("SELECT create_hypertable('dm_time_sdk_log','ts', chunk_time_interval => interval '1 day'    );")
	}
	return &SDKLogRepo{db: db, asyncInsert: stores.NewAsyncInsert[SDK](db, "")}
}
