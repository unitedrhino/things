package relationDB

import "github.com/i-Things/things/shared/stores"

// 角色管理表
type SysTenantRoleInfo struct {
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"` // 角色所属租户编码
	Name       string            `gorm:"column:name;uniqueIndex:tc_ac;type:VARCHAR(100);NOT NULL"`       // 角色名称
	Desc       string            `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`                         //描述

	Status int64               `gorm:"column:status;type:SMALLINT;default:1"` // 状态  1:启用,2:禁用
	Apps   []*SysTenantRoleApp `gorm:"foreignKey:RoleID;references:ID"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;:tc_ac"`
}

func (m *SysTenantRoleInfo) TableName() string {
	return "sys_tenant_role_info"
}

// 应用菜单关联表
type SysTenantRoleApp struct {
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"` // 角色所属租户编码
	RoleID     int64             `gorm:"column:role_id;uniqueIndex:tc_ac;NOT NULL;type:BIGINT"`          // 角色ID
	AppCode    string            `gorm:"column:app_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"`    // 应用编码
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:tc_ac"`
}

func (m *SysTenantRoleApp) TableName() string {
	return "sys_tenant_role_app"
}

// 应用菜单关联表
type SysTenantRoleModule struct {
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"` // 角色所属租户编码
	RoleID     int64             `gorm:"column:role_id;uniqueIndex:tc_ac;NOT NULL;type:BIGINT"`          // 角色ID
	AppCode    string            `gorm:"column:app_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"`    // 应用编码
	ModuleCode string            `gorm:"column:module_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"` // 模块编码
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:tc_ac"`
}

func (m *SysTenantRoleModule) TableName() string {
	return "sys_tenant_role_module"
}

// 应用菜单关联表
type SysTenantRoleMenu struct {
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` // 角色所属租户编码
	RoleID     int64             `gorm:"column:role_id;uniqueIndex:ri_mi;NOT NULL;type:BIGINT"`          // 角色ID
	AppCode    string            `gorm:"column:app_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"`    // 应用编码
	ModuleCode string            `gorm:"column:module_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` // 模块编码
	MenuID     int64             `gorm:"column:menu_id;uniqueIndex:ri_mi;NOT NULL;type:BIGINT"`          // 菜单ID
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:ri_mi"`
}

func (m *SysTenantRoleMenu) TableName() string {
	return "sys_tenant_role_menu"
}

// api权限管理
type SysTenantRoleApi struct {
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` // 角色所属租户编码
	RoleID     int64             `gorm:"column:role_id;uniqueIndex:ri_mi;NOT NULL;type:BIGINT"`          // 角色ID
	AppCode    string            `gorm:"column:app_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"`    // 应用编码
	ModuleCode string            `gorm:"column:module_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` // 模块编码
	ApiID      int64             `gorm:"column:api_id;uniqueIndex:ri_mi;NOT NULL;type:BIGINT"`           // 接口ID
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:ri_mi"`
	Api         *SysTenantAppApi   `gorm:"foreignKey:ID;references:ApiID"`
}

func (m *SysTenantRoleApi) TableName() string {
	return "sys_tenant_role_api"
}
