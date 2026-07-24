// Package usershare 定义跨服务共享的用户设备分享领域事件。
package usershare

// DeviceShareAcceptedDevice 表示一次分享接收成功涉及的设备。
type DeviceShareAcceptedDevice struct {
	ProductID  string `json:"productID"`
	DeviceName string `json:"deviceName"`
}

// DeviceShareAcceptedEvent 表示接收者已成功接收设备分享。
type DeviceShareAcceptedEvent struct {
	EventID         string                      `json:"eventID"`
	ShareToken      string                      `json:"shareToken"`
	SharerUserID    int64                       `json:"sharerUserID"`
	ReceiverUserID  int64                       `json:"receiverUserID"`
	ReceiverAccount string                      `json:"receiverAccount"`
	ProjectID       int64                       `json:"projectID"`
	TenantCode      string                      `json:"tenantCode"`
	UseBy           string                      `json:"useBy"`
	AcceptedAt      int64                       `json:"acceptedAt"`
	Devices         []DeviceShareAcceptedDevice `json:"devices"`
}
