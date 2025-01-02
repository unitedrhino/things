package sdkLogRepo

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/tdengine/sdkLogRepo"
	"github.com/zeromicro/go-zero/core/logx"
)

type SDK struct {
	*deviceLog.SDK
}

func (m *SDK) TableName() string {
	return "dm_sdk_log"
}

type SDKLogRepo struct {
	db          *stores.DB
	asyncInsert *stores.AsyncInsert[SDK]
}

func NewSDKLogRepo(dataSource conf.TSDB) deviceLog.SDKRepo {
	if dataSource.DBType == conf.Tdengine {
		return sdkLogRepo.NewSDKLogRepo(dataSource)
	}
	stores.InitTsConn(dataSource)
	db := stores.GetTsConn(context.Background())
	err := db.AutoMigrate(&SDK{})
	logx.Must(err)
	return &SDKLogRepo{db: db, asyncInsert: stores.NewAsyncInsert[SDK](db, "")}
}
