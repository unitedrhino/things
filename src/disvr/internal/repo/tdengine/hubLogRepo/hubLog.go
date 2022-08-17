package hubLogRepo

import (
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/store"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

type HubLogRepo struct {
	t *clients.Td
	store.HubLogStore
}

func NewHubLogRepo(dataSource string) *HubLogRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &HubLogRepo{t: td}
}
