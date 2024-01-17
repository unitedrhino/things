package relationDB

import (
	"github.com/i-Things/things/shared/stores"
)

// 项目信息表
type SysProjectInfo struct {
	TenantCode  stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);uniqueIndex:pn;NOT NULL"`   // 租户编码
	AdminUserID int64             `gorm:"column:admin_user_id;type:BIGINT;NOT NULL"`                     // 超级管理员id
	ProjectID   stores.ProjectID  `gorm:"column:project_id;type:bigint;NOT NULL"`                        // 项目ID(雪花ID)
	ProjectName string            `gorm:"column:project_name;uniqueIndex:pn;type:varchar(100);NOT NULL"` // 项目名称
	//Region      string            `gorm:"column:region;type:varchar(100);NOT NULL"`      // 项目省市区县
	//Address     string            `gorm:"column:address;type:varchar(512);NOT NULL"`     // 项目详细地址
	Desc string `gorm:"column:desc;type:varchar(100);NOT NULL"` // 项目备注
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:pn"`
}

func (m *SysProjectInfo) TableName() string {
	return "sys_project_info"
}

// 区域信息表
type SysAreaInfo struct {
	TenantCode   stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`   // 租户编码
	AdminUserID  int64             `gorm:"column:admin_user_id;type:BIGINT;NOT NULL"`      // 超级管理员id
	ProjectID    stores.ProjectID  `gorm:"column:project_id;type:bigint;NOT NULL"`         // 所属项目ID(雪花ID)
	AreaID       stores.AreaID     `gorm:"column:area_id;type:bigint;NOT NULL"`            // 区域ID(雪花ID)
	ParentAreaID int64             `gorm:"column:parent_area_id;type:bigint;NOT NULL"`     // 上级区域ID(雪花ID)
	AreaIDPath   string            `gorm:"column:area_id_path;type:varchar(100);NOT NULL"` // 1-2-3-的格式记录顶级区域到当前区域的路径
	AreaName     string            `gorm:"column:area_name;type:varchar(100);NOT NULL"`    // 区域名称
	Position     stores.Point      `gorm:"column:position;type:point;NOT NULL"`            // 区域定位(默认百度坐标系BD09)
	Desc         string            `gorm:"column:desc;type:varchar(100);NOT NULL"`         // 区域备注
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;index"`
	Children    []*SysAreaInfo     `gorm:"foreignKey:ParentAreaID;references:AreaID"`
	Parent      *SysAreaInfo       `gorm:"foreignKey:AreaID;references:ParentAreaID"`
}

func (m *SysAreaInfo) TableName() string {
	return "sys_area_info"
}

// 用户区域权限表
type SysUserArea struct {
	ID         int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` // 租户编码
	UserID     int64             `gorm:"column:user_id;uniqueIndex:ri_mi;type:bigint;NOT NULL"`          // 用户ID(雪花id)
	ProjectID  int64             `gorm:"column:project_id;uniqueIndex:ri_mi;type:bigint;NOT NULL"`       // 所属项目ID(雪花ID)
	AreaID     int64             `gorm:"column:area_id;uniqueIndex:ri_mi;type:bigint;NOT NULL"`          // 区域ID(雪花ID)
	RoleType   int64             `gorm:"column:role_type;type:bigint;NOT NULL"`                          // 角色类型 1 管理员(可以调整本区域及旗下区域的设备区域规划)  2 读写授权(可以对区域下的设备进行操作,但是不能新增或删除)
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;index"`
}

func (m *SysUserArea) TableName() string {
	return "sys_user_area"
}

// 用户项目权限表
type SysUserProject struct {
	ID        int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	UserID    int64 `gorm:"column:user_id;type:bigint;NOT NULL"`    // 用户ID(雪花id)
	ProjectID int64 `gorm:"column:project_id;type:bigint;NOT NULL"` // 所属项目ID(雪花ID)
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;index"`
}

func (m *SysUserProject) TableName() string {
	return "sys_user_project"
}
