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
	return "dm_abnormal_log"
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
	err := db.AutoMigrate(&Abnormal{})
	logx.Must(err)
	return &AbnormalLogRepo{db: db, asyncInsert: stores.NewAsyncInsert[Abnormal](db, "")}
}
