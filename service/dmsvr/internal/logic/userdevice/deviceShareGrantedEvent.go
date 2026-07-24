// 本文件负责构建设备分享授权生效事件，不承载 App 权限修补等消费侧规则。
package userdevicelogic

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strconv"
	"strings"

	"gitee.com/unitedrhino/things/share/domain/usershare"
	"gitee.com/unitedrhino/things/share/topics"
)

// deviceShareGrantedPublisher 约束分享授权生效事件的发布能力，便于独立验证发布行为。
type deviceShareGrantedPublisher interface {
	Publish(ctx context.Context, topic string, arg any) error
}

// buildDeviceShareGrantedEvent 根据本次已生效的分享记录构建领域事件。
func buildDeviceShareGrantedEvent(
	source string,
	shareToken string,
	sharerUserID int64,
	receiverUserID int64,
	receiverAccount string,
	projectID int64,
	tenantCode string,
	useBy string,
	grantedDevices []usershare.DeviceShareGrantedDevice,
	grantedAt int64,
) *usershare.DeviceShareGrantedEvent {
	if len(grantedDevices) == 0 {
		return nil
	}

	devices := append([]usershare.DeviceShareGrantedDevice(nil), grantedDevices...)
	sort.Slice(devices, func(i, j int) bool {
		if devices[i].ProductID == devices[j].ProductID {
			if devices[i].DeviceName == devices[j].DeviceName {
				return devices[i].ShareID < devices[j].ShareID
			}
			return devices[i].DeviceName < devices[j].DeviceName
		}
		return devices[i].ProductID < devices[j].ProductID
	})

	event := &usershare.DeviceShareGrantedEvent{
		Source:          source,
		ShareToken:      shareToken,
		SharerUserID:    sharerUserID,
		ReceiverUserID:  receiverUserID,
		ReceiverAccount: receiverAccount,
		ProjectID:       projectID,
		TenantCode:      tenantCode,
		UseBy:           useBy,
		GrantedAt:       grantedAt,
		Devices:         devices,
	}
	event.EventID = buildDeviceShareGrantedEventID(event)
	return event
}

// buildDeviceShareGrantedEventID 生成与生效时间、设备输入顺序无关的幂等事件编号。
func buildDeviceShareGrantedEventID(event *usershare.DeviceShareGrantedEvent) string {
	parts := []string{
		event.Source,
		event.TenantCode,
		event.ShareToken,
		strconv.FormatInt(event.SharerUserID, 10),
		strconv.FormatInt(event.ReceiverUserID, 10),
		strconv.FormatInt(event.ProjectID, 10),
		event.UseBy,
	}
	for _, device := range event.Devices {
		parts = append(parts, strconv.FormatInt(device.ShareID, 10), device.ProductID, device.DeviceName)
	}
	sum := sha256.Sum256([]byte(strings.Join(parts, "\x00")))
	return "device-share-granted:" + hex.EncodeToString(sum[:])
}

// publishDeviceShareGrantedEvent 发布分享授权生效领域事件。
func publishDeviceShareGrantedEvent(
	ctx context.Context,
	publisher deviceShareGrantedPublisher,
	event *usershare.DeviceShareGrantedEvent,
) error {
	return publisher.Publish(ctx, topics.DmUserDeviceShareGranted, event)
}
