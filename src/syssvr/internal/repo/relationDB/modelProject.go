package relationDB

import (
	"github.com/i-Things/things/shared/stores"
	"time"
)

// 项目信息表
type SysProjectInfo struct {
	TenantCode  stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` // 租户编码
	AdminUserID int64             `gorm:"column:admin_user_id;type:BIGINT;NOT NULL"`                      // 超级管理员id
	ProjectID   stores.ProjectID  `gorm:"column:projectID;type:bigint;NOT NULL"`                          // 项目ID(雪花ID)
	ProjectName string            `gorm:"column:projectName;type:varchar(100);NOT NULL"`                  // 项目名称
	CompanyName string            `gorm:"column:companyName;type:varchar(100);NOT NULL"`                  // 项目所属公司名称
	UserID      int64             `gorm:"column:userID;type:bigint;NOT NULL"`                             // 管理员用户id
	Region      string            `gorm:"column:region;type:varchar(100);NOT NULL"`                       // 项目省市区县
	Address     string            `gorm:"column:address;type:varchar(512);NOT NULL"`                      // 项目详细地址
	Desc        string            `gorm:"column:desc;type:varchar(100);NOT NULL"`                         // 项目备注
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;index"`
}

func (m *SysProjectInfo) TableName() string {
	return "sys_project_info"
}

// 区域信息表
type SysAreaInfo struct {
	TenantCode   stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` // 租户编码
	AdminUserID  int64             `gorm:"column:admin_user_id;type:BIGINT;NOT NULL"`                      // 超级管理员id
	ProjectID    stores.ProjectID  `gorm:"column:projectID;type:bigint;NOT NULL"`                          // 所属项目ID(雪花ID)
	AreaID       stores.AreaID     `gorm:"column:areaID;type:bigint;NOT NULL"`                             // 区域ID(雪花ID)
	ParentAreaID int64             `gorm:"column:parentAreaID;type:bigint;NOT NULL"`                       // 上级区域ID(雪花ID)
	AreaName     string            `gorm:"column:areaName;type:varchar(100);NOT NULL"`                     // 区域名称
	Position     stores.Point      `gorm:"column:position;type:varchar(100);NOT NULL"`                     // 区域定位(默认百度坐标系BD09)
	Desc         string            `gorm:"column:desc;type:varchar(100);NOT NULL"`                         // 区域备注
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
	ID          int64     `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	UserID      int64     `gorm:"column:userID;type:bigint;NOT NULL"`    // 用户ID(雪花id)
	ProjectID   int64     `gorm:"column:projectID;type:bigint;NOT NULL"` // 所属项目ID(雪花ID)
	AreaID      int64     `gorm:"column:areaID;type:bigint;NOT NULL"`    // 区域ID(雪花ID)
	AuthType    int64     `gorm:"column:roleType;type:bigint;NOT NULL"`  // 角色类型 1 管理员  2 读写授权 3 临时授权 4 只读授权
	AuthExpires time.Time `gorm:"column:expires;NOT NULL"`               // 授权过期时间
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;index"`
}

func (m *SysUserArea) TableName() string {
	return "sys_user_auth_area"
}

// 用户项目权限表
type SysUserProject struct {
	ID        int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	UserID    int64 `gorm:"column:userID;type:bigint;NOT NULL"`    // 用户ID(雪花id)
	ProjectID int64 `gorm:"column:projectID;type:bigint;NOT NULL"` // 所属项目ID(雪花ID)
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;index"`
}

func (m *SysUserProject) TableName() string {
	return "sys_user_auth_project"
}
