package hubLogRepo

import (
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/stores"
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
