package deviceDebugLogRepo

import (
	"fmt"
	"github.com/i-Things/things/shared/store/TDengine"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

type DeviceDebugLogRepo struct {
	t *TDengine.Td
}

func NewDeviceDebugLogRepo(dataSource string) *DeviceDebugLogRepo {
	td, err := TDengine.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &DeviceDebugLogRepo{t: td}
}

func getDebugLogStableName(productID string) string {
	return fmt.Sprintf("`model_debug_%s`", productID)
}

func getDebugLogTableName(productID, deviceName string) string {
	return fmt.Sprintf("`device_debug_%s_%s`", productID, deviceName)
}
