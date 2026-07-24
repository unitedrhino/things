// Package usershare 定义跨服务共享的用户设备分享领域事件。
package usershare

const (
	// DeviceShareGrantSourceAccountDirect 表示通过账号直接创建的分享授权。
	DeviceShareGrantSourceAccountDirect = "account_direct"
	// DeviceShareGrantSourceWechatAccept 表示接收者通过微信分享令牌确认的授权。
	DeviceShareGrantSourceWechatAccept = "wechat_accept"
)

// DeviceShareGrantedDevice 表示一次分享授权生效涉及的设备和分享记录。
type DeviceShareGrantedDevice struct {
	ShareID    int64  `json:"shareID"`
	ProductID  string `json:"productID"`
	DeviceName string `json:"deviceName"`
}

// DeviceShareGrantedEvent 表示设备分享授权已对接收者生效。
type DeviceShareGrantedEvent struct {
	EventID         string                     `json:"eventID"`
	Source          string                     `json:"source"`
	ShareToken      string                     `json:"shareToken"`
	SharerUserID    int64                      `json:"sharerUserID"`
	ReceiverUserID  int64                      `json:"receiverUserID"`
	ReceiverAccount string                     `json:"receiverAccount"`
	ProjectID       int64                      `json:"projectID"`
	TenantCode      string                     `json:"tenantCode"`
	UseBy           string                     `json:"useBy"`
	GrantedAt       int64                      `json:"grantedAt"`
	Devices         []DeviceShareGrantedDevice `json:"devices"`
}
