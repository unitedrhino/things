package msgOta

import (
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
)

const (
	TypeReport   = "report"
	TypeUpgrade  = "upgrade"  //固件升级消息下行  返回升级信息，版本、固件地址
	TypeProgress = "progress" //设备端上报升级进度
)

type (
	Req struct {
		deviceMsg.CommonMsg
		Params Params `json:"params,optional"`
	}
	Process struct {
		deviceMsg.CommonMsg
		Params ProcessParams `json:"params,optional"`
	}
	Params struct {
		Version string `json:"version"`
		Module  string `json:"module"`
	}
	ProcessParams struct {
		Step int64  `json:"step"`
		Desc string `json:"desc"`
	}

	//ota下行消息
	Upgrade struct {
		deviceMsg.CommonMsg
		Data *UpgradeData `json:"data,omitempty"`
	}
	UpgradeData struct {
		Version    string  `json:"version"`
		IsDiff     int64   `json:"isDiff"`
		SignMethod string  `json:"signMethod"`
		Extra      string  `json:"extra"`
		Files      []*File `json:"files,omitempty"`
		*File
	}
	File struct {
		Size      int64  `json:"size,omitempty"`
		Name      string `json:"name,omitempty"`
		FileUrl   string `json:"fileUrl,omitempty"`
		FileMd5   string `json:"fileMd5,omitempty"`
		Signature string `json:"signature,omitempty"`
	}
)

func (d *Req) VerifyReqParam() error {

	return nil
}
func (d *Req) GetVersion() string {
	return d.Params.Version
}
func (d *Process) VerifyReqParam() error {
	if d.Params.Step == 0 {
		return errors.Parameter.AddDetail("need add Step")
	}
	return nil
}

// 定义升级包状态常量
const (
	OtaFirmwareStatusNotRequired        = 1 //不需要验证
	OtaFirmwareStatusNotVerified        = 2 //未验证
	OtaFirmwareStatusVerified           = 3 //已验证
	OtaFirmwareStatusVerifying          = 4 //验证中
	OtaFirmwareStatusVerificationFailed = 5 //验证失败
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

// 定义升级批次常量
const (
	ValidateUpgrade = iota + 1 //验证升级包
	BatchUpgrade               //批量升级
)

var JobTypeMap = map[int]string{
	ValidateUpgrade: "验证升级包",
	BatchUpgrade:    "批量升级",
}

// 定义升级任务常量
const (
	UpgradeStatusConfirm = iota + 1
	UpgradeStatusQueued
	UpgradeStatusNotified
	UpgradeStatusInProgress
	UpgradeStatusSucceeded
	UpgradeStatusFailed
	UpgradeStatusCanceled
)

var TaskStatusMap = map[int]string{
	UpgradeStatusConfirm:    "待确认",
	UpgradeStatusQueued:     "待推送",
	UpgradeStatusNotified:   "已推送",
	UpgradeStatusInProgress: "升级中",
	UpgradeStatusSucceeded:  "升级成功",
	UpgradeStatusFailed:     "升级失败",
	UpgradeStatusCanceled:   "已取消",
}

// 定义升级批次常量

/*
静态升级：对于选定的升级范围，仅升级当前满足升级条件的设备。
动态升级：对于选定的升级范围，升级当前满足升级条件的设备，并且持续监测该范围内的设备。只要符合升级条件，物联网平台就会自动推送升级信息。包括但不限于以下设备：
满足升级条件的后续新激活设备。
当前上报的OTA模块版本号不满足升级条件，后续满足升级条件的设备。
*/
type UpgradeType = int64

const (
	StaticUpgrade UpgradeType = iota + 1
	DynamicUpgrade
)

var UpgradeTypeMap = map[int64]string{
	StaticUpgrade:  "静态升级",
	DynamicUpgrade: "动态升级",
}

const (
	AllUpgrade      = iota + 1 //全量升级
	SpecificUpgrade            //定向升级
	GrayUpgrade                //灰度升级
	GroupUpgrade               //分组升级
	AreaUpgrade                //区域升级
)

var UpgradeModeMap = map[int]string{
	AreaUpgrade:     "区域升级",
	AllUpgrade:      "全量升级",
	SpecificUpgrade: "定向升级",
	GrayUpgrade:     "灰度升级",
	GroupUpgrade:    "分组升级",
}

const (
	DiffPackage = iota
	FullPackage
)

var PackageTypeMap = map[int]string{
	FullPackage: "整包",
	DiffPackage: "差包",
}

const ModuleCodeDefault = "default"

type DeviceStatus = int64

const (
	DeviceStatusConfirm    DeviceStatus = iota + 1 //待确认
	DeviceStatusQueued                             //待推送
	DeviceStatusNotified                           //已推送
	DeviceStatusInProgress                         //升级中
	/*
		设备升级完成后，建议立即重启设备，设备上线后，立即上报新的版本号。
		设备上线请求和上报版本请求间隔不能超过2秒。
		重要
			如果设备上报的版本与OTA服务要求的版本一致就认为升级成功，反之认为失败，
			这是物联网平台判断设备升级成功的唯一条件。即使升级进度上报为100%，
			如果不上报新的版本号，可能因为超过设备升级超时时间导致升级失败。
	*/
	DeviceStatusSuccess  //升级成功
	DeviceStatusFailure  //升级失败
	DeviceStatusCanceled //已取消
)

/*
PLANNED：计划中。批次已创建，但是定时时间未到。仅定时静态升级的批次可能返回该值。
IN_PROGRESS：执行中。
COMPLETED：已完成。
CANCELED：已取消。
*/
const (
	JobStatusPlanned    = iota + 1 //计划中。批次已创建，但是定时时间未到。仅定时静态升级的批次可能返回该值。
	JobStatusInProgress            //执行中
	JobStatusCompleted             //已完成
	JobStatusCanceled              //已取消
)
