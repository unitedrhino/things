package msgOta

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg"
)

const (
	TypeReport   = "report"
	TypeUpdate   = "update" //固件升级消息下行  返回升级信息，版本、固件地址
	TypeProgress = "progress"
)

type (
	Req struct {
		deviceMsg.CommonMsg
		Params params `json:"params,optional"`
	}
	Process struct {
		deviceMsg.CommonMsg
		Params processParams `json:"params,optional"`
	}
	params struct {
		ID      int64  `json:"id"`
		Version string `json:"version"`
		Module  string `json:"module"`
	}
	processParams struct {
		ID     int64  `json:"id"`
		Step   int64  `json:"step"`
		Desc   string `json:"desc"`
		Module string `json:"module"`
	}
)

func (d *Req) VerifyReqParam() error {
	if d.Params.Module == "" {
		return errors.Parameter.AddDetail("need add module")
	}

	return nil
}
func (d *Req) GetVersion() string {
	return d.Params.Version
}
func (d *Process) VerifyReqParam() error {
	if d.Params.Module == "" {
		return errors.Parameter.AddDetail("need add module")
	}
	if d.Params.Step == 0 {
		return errors.Parameter.AddDetail("need add Step")
	}
	return nil
}

// 定义升级包状态常量
const (
	OtaFirmwareStatusNotRequired        = -1
	OtaFirmwareStatusNotVerified        = 0
	OtaFirmwareStatusVerified           = 1
	OtaFirmwareStatusVerifying          = 2
	OtaFirmwareStatusVerificationFailed = 3
)

// 定义升级包状态映射
var OtaFirmwareStatusMap = map[int]string{
	OtaFirmwareStatusNotRequired:        "不需要验证",
	OtaFirmwareStatusNotVerified:        "未验证",
	OtaFirmwareStatusVerified:           "已验证",
	OtaFirmwareStatusVerifying:          "验证中",
	OtaFirmwareStatusVerificationFailed: "验证失败",
}

// 根据状态值返回中文字符串
func GetOtaFirmwareStatusString(status int) string {
	if statusString, ok := OtaFirmwareStatusMap[status]; ok {
		return statusString
	}
	return "未知状态"
}
