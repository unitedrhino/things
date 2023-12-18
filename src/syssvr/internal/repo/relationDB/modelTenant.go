package relationDB

import "github.com/i-Things/things/shared/stores"

// 租户信息表
type SysTenantInfo struct {
	ID          int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`   // id编号
	Code        string `gorm:"column:code;uniqueIndex;type:VARCHAR(100);NOT NULL"` // 租户编码
	Name        string `gorm:"column:name;uniqueIndex;type:VARCHAR(100);NOT NULL"` // 租户名称
	AdminUserID int64  `gorm:"column:admin_user_id;type:BIGINT;NOT NULL"`          // 超级管理员id
	BaseUrl     string `gorm:"column:base_url;type:VARCHAR(100);NOT NULL"`         //应用首页
	LogoUrl     string `gorm:"column:logo_url;type:VARCHAR(100);NOT NULL"`         //应用logo地址
	Desc        string `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`             //应用描述
	stores.Time
}

func (m *SysTenantInfo) TableName() string {
	return "sys_tenant_info"
}

// 租户下的应用列表
type SysTenantApp struct {
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"` // 租户编码
	AppCode    string            `gorm:"column:app_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"`    // 应用编码 这里只关联主应用,主应用授权,子应用也授权了
	stores.Time
}

func (m *SysTenantApp) TableName() string {
	return "sys_tenant_app"
}

// 租户下的邮箱配置
type SysTenantConfig struct {
	ID             int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`         // id编号
	TenantCode     stores.TenantCode `gorm:"column:tenant_code;uniqueIndex;type:VARCHAR(50);NOT NULL"` // 租户编码
	Email          *SysTenantEmail   `gorm:"embedded;embeddedPrefix:email_"`
	RegisterRoleID int64             `gorm:"column:register_role_id;type:BIGINT;NOT NULL"` //注册分配的角色id
	stores.Time
}

type SysTenantEmail struct {
	From     string `gorm:"column:from;type:VARCHAR(50);NOT NULL"`     // 发件人  你自己要发邮件的邮箱
	Host     string `gorm:"column:host;type:VARCHAR(50);NOT NULL"`     // 服务器地址 例如 smtp.qq.com  请前往QQ或者你要发邮件的邮箱查看其smtp协议
	Secret   string `gorm:"column:secret;type:VARCHAR(50);NOT NULL"`   // 密钥    用于登录的密钥 最好不要用邮箱密码 去邮箱smtp申请一个用于登录的密钥
	Nickname string `gorm:"column:nickname;type:VARCHAR(50);NOT NULL"` // 昵称    发件人昵称 通常为自己的邮箱
	Port     int64  `gorm:"column:port;type:int;default:465"`          // 端口     请前往QQ或者你要发邮件的邮箱查看其smtp协议 大多为 465
	IsSSL    int64  `gorm:"column:is_ssl;type:int;default:2"`          // 是否SSL   是否开启SSL
}

func (m *SysTenantConfig) TableName() string {
	return "sys_tenant_config"
}
