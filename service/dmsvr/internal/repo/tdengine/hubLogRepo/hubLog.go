package hubLogRepo

import (
	"fmt"
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"sync"
)

type HubLogRepo struct {
	t    *clients.Td
	once sync.Once
}

func NewHubLogRepo(dataSource conf.TSDB) *HubLogRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &HubLogRepo{t: td}
}

func (h *HubLogRepo) GetLogStableName() string {
	return fmt.Sprintf("`model_common_hublog`")
}

func (h *HubLogRepo) GetLogTableName(productID, deviceName string) string {
	return fmt.Sprintf("`device_hublog_%s_%s`", productID, deviceName)
}
