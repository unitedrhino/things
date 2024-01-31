package relationDB

import "gitee.com/i-Things/core/shared/stores"

// 示例
type UdExample struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
}

// 设备维护工单 device Maintenance Work Order
type UdOpsWorkOrder struct {
	ID          int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	TenantCode  stores.TenantCode `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL"`         // 租户编码
	RaiseUserID int64             `gorm:"column:raise_user_id;type:BIGINT;NOT NULL"`                  // 问题提出的用户
	ProjectID   stores.ProjectID  `gorm:"column:project_id;type:bigint;default:0;NOT NULL"`           // 项目ID(雪花ID)
	AreaID      stores.AreaID     `gorm:"column:area_id;type:bigint;default:0;NOT NULL"`              // 项目区域ID(雪花ID)
	Number      string            `gorm:"column:number;uniqueIndex:number;type:VARCHAR(50);NOT NULL"` //编号
	Params      string            `gorm:"column:params;type:varchar(1024);NOT NULL"`                  // 参数 json格式
	Type        string            `gorm:"column:type;type:varchar(100);NOT NULL"`                     // 工单类型: deviceMaintenance:设备维修工单
	IssueDesc   string            `gorm:"column:issue_desc;type:varchar(2000);NOT NULL"`
	Status      int64             `gorm:"column:status;type:BIGINT;default:1"` //状态 1:待处理 2:处理中 3:已完成
	stores.SoftTime
}

func (m *UdOpsWorkOrder) TableName() string {
	return "ud_ops_work_order"
}

type UdUserCollectDevice struct {
	ID         int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	TenantCode stores.TenantCode `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL"`                              // 租户编码
	ProjectID  stores.ProjectID  `gorm:"column:project_id;type:bigint;default:0;NOT NULL"`                                // 项目ID(雪花ID)
	UserID     int64             `gorm:"column:user_id;type:BIGINT;uniqueIndex:product_id_deviceName;NOT NULL"`           // 问题提出的用户
	ProductID  string            `gorm:"column:product_id;type:char(11);uniqueIndex:product_id_deviceName;NOT NULL"`      // 产品id
	DeviceName string            `gorm:"column:device_name;uniqueIndex:product_id_deviceName;type:varchar(100);NOT NULL"` // 设备名称
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:product_id_deviceName"`
}

func (m *UdUserCollectDevice) TableName() string {
	return "ud_user_collect_device"
}
