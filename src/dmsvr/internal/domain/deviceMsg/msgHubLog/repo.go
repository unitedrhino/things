// Package device 设备操作日志
package msgHubLog

import (
	"context"
	"gitee.com/i-Things/core/shared/def"
	"time"
)

type (
	HubLog struct {
		ProductID  string    `json:"productID"`  // 产品id
		DeviceName string    `json:"deviceName"` // 设备名称
		Content    string    `json:"content"`    // 具体信息
		Topic      string    `json:"topic"`      // 主题
		Action     string    `json:"action"`     // 操作类型
		Timestamp  time.Time `json:"timestamp"`  // 操作时间
		RequestID  string    `json:"requestID"`  // 请求ID
		TranceID   string    `json:"tranceID"`   // 服务器端事务id
		ResultType int64     `json:"resultType"` // 请求结果状态,200为成功
	}
	HubFilter struct {
		ProductID  string   // 产品id
		DeviceName string   // 设备名称
		Actions    []string //过滤操作类型 connected:上线 disconnected:下线  property:属性 event:事件 action:操作 thing:物模型提交的操作为匹配的日志
		Topics     []string //过滤主题
		Content    string   //过滤内容
		RequestID  string   //过滤请求ID
	}
	HubLogRepo interface {
		InitProduct(ctx context.Context, productID string) error
		DropProduct(ctx context.Context, productID string) error
		DropDevice(ctx context.Context, productID string, deviceName string) error
		GetDeviceLog(ctx context.Context, filter HubFilter, page def.PageInfo2) ([]*HubLog, error)
		GetCountLog(ctx context.Context, filter HubFilter, page def.PageInfo2) (int64, error)
		Insert(ctx context.Context, data *HubLog) error
	}
)
