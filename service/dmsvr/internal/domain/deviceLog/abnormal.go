// Package device 设备操作日志
package deviceLog

import (
	"context"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/things/share/devices"
	"time"
)

type (
	Abnormal struct {
		TenantCode  dataType.TenantCode    `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL"`                        // 租户编码
		ProjectID   dataType.ProjectID     `gorm:"column:project_id;index:project_id_area_id;type:bigint;default:0;NOT NULL"` // 项目ID(雪花ID)
		AreaID      dataType.AreaID        `gorm:"column:area_id;index:project_id_area_id;type:bigint;default:0;NOT NULL"`    // 项目区域ID(雪花ID)
		AreaIDPath  dataType.AreaIDPath    `gorm:"column:area_id_path;type:varchar(100);default:'';NOT NULL"`                 // 项目区域ID路径(雪花ID)
		BelongGroup map[string]def.IDsInfo `gorm:"column:belong_group;type:json;serializer:json;default:'{}'"`
		ProductID   string                 `gorm:"column:product_id;type:varchar(100);NOT NULL" json:"productID,omitempty"`   // 产品id
		DeviceName  string                 `gorm:"column:device_name;type:varchar(100);NOT NULL" json:"deviceName,omitempty"` // 设备名称
		Action      def.Bool               `gorm:"column:user_id;type:BIGINT;default:1" json:"action,omitempty"`              //触发1还是解除2
		Type        string                 `gorm:"column:type;type:varchar(100);NOT NULL" json:"type,omitempty"`              // 异常类型
		Timestamp   time.Time              `gorm:"column:ts;NOT NULL;" json:"timestamp"`                                      // 操作时间
		TraceID     string                 `gorm:"column:trace_id;type:varchar(100);NOT NULL" json:"traceID,omitempty"`       // 服务器端事务id
		Reason      string                 `gorm:"column:reason;type:varchar(100);NOT NULL" json:"reason,omitempty"`          //原因

	}
	AbnormalFilter struct {
		TenantCode  string
		ProjectID   int64   `json:"projectID,omitempty"`
		AreaID      int64   `json:"areaID,omitempty"`
		AreaIDPath  string  `json:"areaIDPath,omitempty"`
		AreaIDs     []int64 `json:"areaIDs"`
		BelongGroup map[string]def.IDsInfo
		ProductID   string // 产品id
		ProductIDs  []string
		DeviceName  string // 设备名称
		Action      int64
		Type        string `json:"type,omitempty"` // 异常类型
		Reason      string
	}

	AbnormalRepo interface {
		GetDeviceLog(ctx context.Context, filter AbnormalFilter, page def.PageInfo2) ([]*Abnormal, error)
		GetCountLog(ctx context.Context, filter AbnormalFilter, page def.PageInfo2) (int64, error)
		Insert(ctx context.Context, data *Abnormal) error
		ManageRepo
		UpdateDevice(ctx context.Context, devices []*devices.Core, affiliation devices.Affiliation) error
		UpdateDevices(ctx context.Context, devices []*devices.Info) error
		VersionUpdate(ctx context.Context, version string) error
		//ModifyRepo
	}
)
