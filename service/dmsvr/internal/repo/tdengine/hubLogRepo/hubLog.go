package hubLogRepo

import (
	"gitee.com/i-Things/core/shared/clients"
	"gitee.com/i-Things/core/shared/conf"
	"gitee.com/i-Things/core/shared/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

type HubLogRepo struct {
	t *clients.Td
	stores.HubLogStore
}

func NewHubLogRepo(dataSource conf.TSDB) *HubLogRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &HubLogRepo{t: td}
}
