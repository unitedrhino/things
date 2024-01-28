package relationDB

import "github.com/i-Things/things/shared/stores"

// 功能权限范围
type SysAccessInfo struct {
	ID         int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`             // 编号
	Name       string `gorm:"column:name;type:VARCHAR(100);NOT NULL"`                       // 请求名称
	Code       string `gorm:"column:code;type:VARCHAR(100);uniqueIndex:app_route;NOT NULL"` // 请求名称
	Group      string `gorm:"column:group;type:VARCHAR(100);NOT NULL"`                      // 接口组
	IsNeedAuth int64  `gorm:"column:is_need_auth;type:BIGINT;default:1;NOT NULL"`           // 是否需要认证（1是 2否）
	Desc       string `gorm:"column:desc;type:VARCHAR(500);NOT NULL"`                       // 备注
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:app_route"`
	Apis        []*SysApiInfo      `gorm:"foreignKey:AccessCode;references:Code"`
}

func (m *SysAccessInfo) TableName() string {
	return "sys_access_info"
}

// 接口管理
type SysApiInfo struct {
	ID           int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`          // 编号
	AccessCode   string `gorm:"column:api_scope_code;type:VARCHAR(50);NOT NULL"`           // 范围编码
	Method       string `gorm:"column:method;uniqueIndex:route;type:VARCHAR(50);NOT NULL"` // 请求方式（1 GET 2 POST 3 HEAD 4 OPTIONS 5 PUT 6 DELETE 7 TRACE 8 CONNECT 9 其它）
	Route        string `gorm:"column:route;uniqueIndex:route;type:VARCHAR(100);NOT NULL"` // 路由
	Name         string `gorm:"column:name;type:VARCHAR(100);NOT NULL"`                    // 请求名称
	BusinessType int64  `gorm:"column:business_type;type:BIGINT;NOT NULL"`                 // 业务类型（1新增 2修改 3删除 4查询 5其它）
	Desc         string `gorm:"column:desc;type:VARCHAR(500);NOT NULL"`                    // 备注
	IsAuthTenant int64  `gorm:"column:is_auth_tenant;type:BIGINT;default:1;NOT NULL"`      // 是否可以授权给普通租户
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:route"`
}

func (m *SysApiInfo) TableName() string {
	return "sys_api_info"
}

// 应用菜单关联表
type SysTenantAccess struct {
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                         // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:tenant_scope;type:VARCHAR(50);NOT NULL;"`   // 租户编码
	AccessCode string            `gorm:"column:api_scope_code;uniqueIndex:tenant_scope;type:VARCHAR(50);NOT NULL"` // 范围编码
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:tenant_scope"`
}

func (m *SysTenantAccess) TableName() string {
	return "sys_tenant_access"
}
