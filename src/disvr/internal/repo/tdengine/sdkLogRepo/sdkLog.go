package sdkLogRepo

import (
	"fmt"
	"github.com/i-Things/things/shared/clients"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

type SDKLogRepo struct {
	t *clients.Td
}

func NewSDKLogRepo(dataSource string) *SDKLogRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &SDKLogRepo{t: td}
}

func getSDKLogStableName(productID string) string {
	return fmt.Sprintf("`model_sdklog_%s`", productID)
}

func getSDKLogTableName(productID, deviceName string) string {
	return fmt.Sprintf("`sdk_log_%s_%s`", productID, deviceName)
}
