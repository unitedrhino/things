package hubLogRepo

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/tdengine/hubLogRepo"
	"github.com/zeromicro/go-zero/core/logx"
)

type Hub struct {
	*deviceLog.Hub
}

func (m *Hub) TableName() string {
	return "dm_hub_log"
}

type HubLogRepo struct {
	db          *stores.DB
	asyncInsert *stores.AsyncInsert[Hub]
}

func NewHubLogRepo(dataSource conf.TSDB) deviceLog.HubRepo {
	if dataSource.DBType == conf.Tdengine {
		return hubLogRepo.NewHubLogRepo(dataSource)
	}
	stores.InitTsConn(dataSource)
	db := stores.GetTsConn(context.Background())
	err := db.AutoMigrate(&Hub{})
	logx.Must(err)
	return &HubLogRepo{db: db, asyncInsert: stores.NewAsyncInsert[Hub](db, "")}
}

func (h *HubLogRepo) GetLogStableName() string {
	return fmt.Sprintf("`model_common_hublog`")
}

func (h *HubLogRepo) GetLogTableName(productID, deviceName string) string {
	return fmt.Sprintf("`device_hublog_%s_%s`", productID, deviceName)
}
