// Package device 这个文件提供基础信息的物模型定义
package baseParam

// BasicParam 小程序或 App 展示设备详细信息时，一般会展示设备的 MAC 地址、IMEI 号、时区等基础信息。设备信息上报使用的 Topic：
//上行请求 Topic： $thing/up/property/{ProductID}/{DeviceName}
//下行响应 Topic： $thing/down/property/{ProductID}/{DeviceName}
type BasicParam struct {
	Name           string            `json:"name"`           //设备名(是否保留待定)
	Imei           string            `json:"imei"`           //设备的 IMEI 号信息，非必填项
	FwVer          string            `json:"fwVer"`          //mcu固件版本
	ModuleHardInfo string            `json:"moduleHardInfo"` //模组具体硬件型号
	ModuleSoftInfo string            `json:"moduleSoftInfo"` //模组软件版本
	Mac            string            `json:"mac"`            //设备的 MAC 信息，非必填项
	DeviceLabel    map[string]string `json:"deviceLabel"`    //设备商自定义的产品基础信息，以 KV 方式上报
}

//todo 具体实现参考设备基础信息上报,https://cloud.tencent.com/document/product/1081/34916
