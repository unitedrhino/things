package hubLogRepo

import (
	"fmt"
	"github.com/i-Things/things/shared/clients"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

type HubLogRepo struct {
	t *clients.Td
}

func NewHubLogRepo(dataSource string) *HubLogRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &HubLogRepo{t: td}
}

func getLogStableName(productID string) string {
	return fmt.Sprintf("`model_hublog_%s`", productID)
}

func getLogTableName(productID, deviceName string) string {
	return fmt.Sprintf("`hub_log_%s_%s`", productID, deviceName)
}
