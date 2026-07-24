package usershare

import (
	"encoding/json"
	"testing"
)

// TestDeviceShareGrantedEventJSONContract 验证跨服务事件的 JSON 字段契约。
func TestDeviceShareGrantedEventJSONContract(t *testing.T) {
	event := DeviceShareGrantedEvent{
		EventID:         "event-1",
		Source:          DeviceShareGrantSourceWechatAccept,
		ShareToken:      "share-token-1",
		SharerUserID:    101,
		ReceiverUserID:  202,
		ReceiverAccount: "receiver",
		ProjectID:       303,
		TenantCode:      "default",
		UseBy:           "wechat_single_device",
		GrantedAt:       1717500000,
		Devices: []DeviceShareGrantedDevice{
			{ShareID: 404, ProductID: "product-1", DeviceName: "device-1"},
		},
	}

	body, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	want := `{"eventID":"event-1","source":"wechat_accept","shareToken":"share-token-1","sharerUserID":101,"receiverUserID":202,"receiverAccount":"receiver","projectID":303,"tenantCode":"default","useBy":"wechat_single_device","grantedAt":1717500000,"devices":[{"shareID":404,"productID":"product-1","deviceName":"device-1"}]}`
	if string(body) != want {
		t.Fatalf("json.Marshal() = %s, want %s", body, want)
	}
}
