package usershare

import (
	"encoding/json"
	"testing"
)

func TestDeviceShareAcceptedEventJSONContract(t *testing.T) {
	event := DeviceShareAcceptedEvent{
		EventID:         "event-1",
		ShareToken:      "share-token-1",
		SharerUserID:    101,
		ReceiverUserID:  202,
		ReceiverAccount: "receiver",
		ProjectID:       303,
		TenantCode:      "default",
		UseBy:           "wechat_single_device",
		AcceptedAt:      1717500000,
		Devices: []DeviceShareAcceptedDevice{
			{ProductID: "product-1", DeviceName: "device-1"},
		},
	}

	body, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	want := `{"eventID":"event-1","shareToken":"share-token-1","sharerUserID":101,"receiverUserID":202,"receiverAccount":"receiver","projectID":303,"tenantCode":"default","useBy":"wechat_single_device","acceptedAt":1717500000,"devices":[{"productID":"product-1","deviceName":"device-1"}]}`
	if string(body) != want {
		t.Fatalf("json.Marshal() = %s, want %s", body, want)
	}
}
