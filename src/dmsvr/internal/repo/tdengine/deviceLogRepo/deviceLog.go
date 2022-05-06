package deviceLogRepo

import (
	"fmt"
	"github.com/i-Things/things/shared/store/TDengine"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

type DeviceLogRepo struct {
	t *TDengine.Td
}

func NewDeviceLogRepo(dataSource string) *DeviceLogRepo {
	td, err := TDengine.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &DeviceLogRepo{t: td}
}

func getLogStableName(productID string) string {
	return fmt.Sprintf("`model_log_%s`", productID)
}

func getLogTableName(productID, deviceName string) string {
	return fmt.Sprintf("`device_log_%s_%s`", productID, deviceName)
}
