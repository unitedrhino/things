package deviceMsgEvent

import (
	"gitee.com/i-Things/share/domain/deviceMsg/msgThing"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"strings"
	"testing"
)

func TestThingLogic_HandlePropertyReportInfo(t *testing.T) {
	req := map[string]any{
		"position": map[string]interface{}{
			"coordinateSystem": "WGS84",
			"longitude":        "12",
			"latitude":         "34",
		},
	}
	diDeviceBasicInfoDo := &msgThing.DeviceBasicInfo{}
	if err := mapstructure.Decode(req, diDeviceBasicInfoDo); err != nil {
		if strings.Contains(err.Error(), "expected type") && req["position"] != nil {
			pos, ok := req["position"].(map[string]interface{})
			if !ok {
				return
			}
			pos["latitude"] = cast.ToFloat64(pos["latitude"])
			pos["longitude"] = cast.ToFloat64(pos["longitude"])
			req["position"] = pos
			err = mapstructure.Decode(req, diDeviceBasicInfoDo)
			if err != nil {
				return
			}
		}
	}
	return
}
