package userdevicelogic

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strconv"
	"strings"

	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/domain/usershare"
	"gitee.com/unitedrhino/things/share/topics"
)

// deviceShareAcceptedPublisher 约束分享接收成功事件的发布能力，便于独立验证发布行为。
type deviceShareAcceptedPublisher interface {
	Publish(ctx context.Context, topic string, arg any) error
}

// buildDeviceShareAcceptedEvent 根据本次实际接收成功的设备构建领域事件。
func buildDeviceShareAcceptedEvent(
	in *dm.UserDeviceShareMultiAcceptReq,
	info *dm.UserDeviceShareMultiInfo,
	tenantCode string,
	acceptedDevices []*dm.DeviceShareInfo,
	acceptedAt int64,
) *usershare.DeviceShareAcceptedEvent {
	if in == nil || info == nil || len(acceptedDevices) == 0 {
		return nil
	}

	devices := make([]usershare.DeviceShareAcceptedDevice, 0, len(acceptedDevices))
	for _, device := range acceptedDevices {
		if device == nil {
			continue
		}
		devices = append(devices, usershare.DeviceShareAcceptedDevice{
			ProductID:  device.ProductID,
			DeviceName: device.DeviceName,
		})
	}
	if len(devices) == 0 {
		return nil
	}
	sort.Slice(devices, func(i, j int) bool {
		if devices[i].ProductID == devices[j].ProductID {
			return devices[i].DeviceName < devices[j].DeviceName
		}
		return devices[i].ProductID < devices[j].ProductID
	})

	event := &usershare.DeviceShareAcceptedEvent{
		ShareToken:      in.ShareToken,
		SharerUserID:    info.UserID,
		ReceiverUserID:  in.SharedUserID,
		ReceiverAccount: in.SharedUserAccount,
		ProjectID:       info.ProjectID,
		TenantCode:      tenantCode,
		UseBy:           info.UseBy,
		AcceptedAt:      acceptedAt,
		Devices:         devices,
	}
	event.EventID = buildDeviceShareAcceptedEventID(event)
	return event
}

// buildDeviceShareAcceptedEventID 生成与接收时间、设备输入顺序无关的幂等事件编号。
func buildDeviceShareAcceptedEventID(event *usershare.DeviceShareAcceptedEvent) string {
	parts := []string{
		event.TenantCode,
		event.ShareToken,
		strconv.FormatInt(event.SharerUserID, 10),
		strconv.FormatInt(event.ReceiverUserID, 10),
		strconv.FormatInt(event.ProjectID, 10),
		event.UseBy,
	}
	for _, device := range event.Devices {
		parts = append(parts, device.ProductID, device.DeviceName)
	}
	sum := sha256.Sum256([]byte(strings.Join(parts, "\x00")))
	return "device-share-accepted:" + hex.EncodeToString(sum[:])
}

// publishDeviceShareAcceptedEvent 发布分享接收成功领域事件。
func publishDeviceShareAcceptedEvent(
	ctx context.Context,
	publisher deviceShareAcceptedPublisher,
	event *usershare.DeviceShareAcceptedEvent,
) error {
	return publisher.Publish(ctx, topics.DmUserDeviceShareAccepted, event)
}
