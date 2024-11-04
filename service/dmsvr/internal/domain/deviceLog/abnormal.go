// Package device 设备操作日志
package deviceLog

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"time"
)

type (
	Abnormal struct {
		ProductID  string    `json:"productID,omitempty"`  // 产品id
		DeviceName string    `json:"deviceName,omitempty"` // 设备名称
		Action     bool      `json:"action,omitempty"`     //触发1还是解除2
		Type       string    `json:"type,omitempty"`       // 异常类型
		Timestamp  time.Time `json:"timestamp"`            // 操作时间
		TraceID    string    `json:"traceID,omitempty"`    // 服务器端事务id
		Reason     string    `json:"reason,omitempty"`     //原因
	}
	AbnormalFilter struct {
		ProductID  string // 产品id
		DeviceName string // 设备名称
		Action     bool
		Type       string `json:"type,omitempty"` // 异常类型
		Reason     string
	}

	AbnormalRepo interface {
		GetDeviceLog(ctx context.Context, filter AbnormalFilter, page def.PageInfo2) ([]*Abnormal, error)
		GetCountLog(ctx context.Context, filter AbnormalFilter, page def.PageInfo2) (int64, error)
		Insert(ctx context.Context, data *Abnormal) error
		ManageRepo
		//ModifyRepo
	}
)
