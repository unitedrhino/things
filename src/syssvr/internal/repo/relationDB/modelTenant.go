package relationDB

import (
	"database/sql"
	"github.com/i-Things/things/shared/stores"
	"time"
)

// 租户信息表
type SysTenantInfo struct {
	ID          int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`        // id编号
	Code        string `gorm:"column:code;uniqueIndex:code;type:VARCHAR(100);NOT NULL"` // 租户编码
	Name        string `gorm:"column:name;uniqueIndex:name;type:VARCHAR(100);NOT NULL"` // 租户名称
	AdminUserID int64  `gorm:"column:admin_user_id;type:BIGINT;NOT NULL"`               // 超级管理员id
	BaseUrl     string `gorm:"column:base_url;type:VARCHAR(100);NOT NULL"`              //应用首页
	LogoUrl     string `gorm:"column:logo_url;type:VARCHAR(100);NOT NULL"`              //应用logo地址
	Desc        string `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`                  //应用描述
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:code;uniqueIndex:name"`
}

func (m *SysTenantInfo) TableName() string {
	return "sys_tenant_info"
}

// 租户下的应用列表
type SysTenantApp struct {
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"` // 租户编码
	AppCode    string            `gorm:"column:app_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"`    // 应用编码 这里只关联主应用,主应用授权,子应用也授权了
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:ri_mi"`
}

func (m *SysTenantApp) TableName() string {
	return "sys_tenant_app"
}

// 租户下的应用列表
type SysTenantAppModule struct {
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"` // 租户编码
	SysAppModule
}

func (m *SysTenantAppModule) TableName() string {
	return "sys_tenant_app_module"
}

// 接口管理
type SysTenantAppApi struct {
	TempLateID int64             `gorm:"column:template_id;type:BIGINT;NOT NULL"`                            // 模板id
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:app_route;type:VARCHAR(50);NOT NULL"` // 租户编码
	AppCode    string            `gorm:"column:app_code;uniqueIndex:app_route;type:VARCHAR(50);NOT NULL"`    // 应用编码 这里只关联主应用,主应用授权,子应用也授权了
	SysModuleApi
}

func (m *SysTenantAppApi) TableName() string {
	return "sys_tenant_app_api"
}

// 菜单管理表
type SysTenantAppMenu struct {
	TempLateID int64             `gorm:"column:template_id;type:BIGINT;NOT NULL"`      // 模板id
	TenantCode stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"` // 租户编码
	AppCode    string            `gorm:"column:app_code;type:VARCHAR(50);NOT NULL"`    // 应用编码 这里只关联主应用,主应用授权,子应用也授权了
	SysModuleMenu
}

func (m *SysTenantAppMenu) TableName() string {
	return "sys_tenant_app_menu"
}

// 用户登录信息表
type SysTenantUserInfo struct {
	TenantCode    stores.TenantCode    `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL;uniqueIndex:tc_un;uniqueIndex:tc_email;uniqueIndex:tc_phone;uniqueIndex:tc_wui;uniqueIndex:tc_woi"` // 租户编码
	UserID        int64                `gorm:"column:user_id;primary_key;AUTO_INCREMENT;type:BIGINT;NOT NULL"`                                                                                 // 用户id
	UserName      sql.NullString       `gorm:"column:user_name;uniqueIndex:tc_un;type:VARCHAR(20)"`                                                                                            // 登录用户名
	NickName      string               `gorm:"column:nick_name;type:VARCHAR(60);NOT NULL"`                                                                                                     // 用户的昵称
	Password      string               `gorm:"column:password;type:CHAR(32);NOT NULL"`                                                                                                         // 登录密码
	Email         sql.NullString       `gorm:"column:email;uniqueIndex:tc_email;type:VARCHAR(255)"`                                                                                            // 邮箱
	Phone         sql.NullString       `gorm:"column:phone;uniqueIndex:tc_phone;type:VARCHAR(20)"`                                                                                             // 手机号
	WechatUnionID sql.NullString       `gorm:"column:wechat_union_id;uniqueIndex:tc_wui;type:VARCHAR(20)"`                                                                                     // 微信union id
	WechatOpenID  sql.NullString       `gorm:"column:wechat_open_id;uniqueIndex:tc_woi;type:VARCHAR(20)"`                                                                                      // 微信union id
	LastIP        string               `gorm:"column:last_ip;type:VARCHAR(40);NOT NULL"`                                                                                                       // 最后登录ip
	RegIP         string               `gorm:"column:reg_ip;type:VARCHAR(40);NOT NULL"`                                                                                                        // 注册ip
	Sex           int64                `gorm:"column:sex;type:SMALLINT;default:3;NOT NULL"`                                                                                                    // 用户的性别，值为1时是男性，值为2时是女性，其他值为未知
	City          string               `gorm:"column:city;type:VARCHAR(50);NOT NULL"`                                                                                                          // 用户所在城市
	Country       string               `gorm:"column:country;type:VARCHAR(50);NOT NULL"`                                                                                                       // 用户所在国家
	Province      string               `gorm:"column:province;type:VARCHAR(50);NOT NULL"`                                                                                                      // 用户所在省份
	Language      string               `gorm:"column:language;type:VARCHAR(50);NOT NULL"`                                                                                                      // 用户的语言，简体中文为zh_CN
	HeadImg       string               `gorm:"column:head_img;type:VARCHAR(256);NOT NULL"`                                                                                                     // 用户头像
	Role          int64                `gorm:"column:role;type:BIGINT;NOT NULL"`                                                                                                               // 用户默认角色（默认使用该角色）
	IsAllData     int64                `gorm:"column:is_all_data;type:SMALLINT;default:1;NOT NULL"`                                                                                            // 是否所有数据权限（1是，2否）
	Roles         []*SysTenantUserRole `gorm:"foreignKey:UserID;references:UserID"`
	Tenant        *SysTenantInfo       `gorm:"foreignKey:Code;references:TenantCode"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:tc_un;uniqueIndex:tc_email;uniqueIndex:tc_phone;uniqueIndex:tc_wui;uniqueIndex:tc_woi"`
}

func (m *SysTenantUserInfo) TableName() string {
	return "sys_tenant_user_info"
}

// 应用菜单关联表
type SysTenantUserRole struct {
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`      // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL;"`         // 租户编码
	UserID     int64             `gorm:"column:user_id;uniqueIndex:ri_mi;NOT NULL;type:BIGINT"` // 用户ID
	RoleID     int64             `gorm:"column:role_id;uniqueIndex:ri_mi;NOT NULL;type:BIGINT"` // 角色ID
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:ri_mi"`
}

func (m *SysTenantUserRole) TableName() string {
	return "sys_tenant_user_role"
}

// 登录日志管理
type SysTenantLoginLog struct {
	ID            int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`       // 编号
	TenantCode    stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`           // 租户编码
	AppCode       string            `gorm:"column:app_code;NOT NULL;type:VARCHAR(50)"`              // 应用ID
	UserID        int64             `gorm:"column:user_id;type:BIGINT;NOT NULL"`                    // 用户id
	UserName      string            `gorm:"column:user_name;type:VARCHAR(50)"`                      // 登录账号
	IpAddr        string            `gorm:"column:ip_addr;type:VARCHAR(50)"`                        // 登录IP地址
	LoginLocation string            `gorm:"column:login_location;type:VARCHAR(100)"`                // 登录地点
	Browser       string            `gorm:"column:browser;type:VARCHAR(50)"`                        // 浏览器类型
	Os            string            `gorm:"column:os;type:VARCHAR(50)"`                             // 操作系统
	Code          int64             `gorm:"column:code;type:BIGINT;default:200;NOT NULL"`           // 登录状态（200成功 其它失败）
	Msg           string            `gorm:"column:msg;type:VARCHAR(255)"`                           // 提示消息
	CreatedTime   time.Time         `gorm:"column:created_time;default:CURRENT_TIMESTAMP;NOT NULL"` // 登录时间
}

func (m *SysTenantLoginLog) TableName() string {
	return "sys_tenant_login_log"
}

// 操作日志管理
type SysTenantOperLog struct {
	ID           int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`       // 编号
	TenantCode   stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`           // 租户编码
	AppCode      string            `gorm:"column:app_code;NOT NULL;type:VARCHAR(50)"`              // 应用ID
	OperUserID   int64             `gorm:"column:oper_user_id;type:BIGINT;NOT NULL"`               // 用户id
	OperUserName string            `gorm:"column:oper_user_name;type:VARCHAR(50)"`                 // 操作人员名称
	OperName     string            `gorm:"column:oper_name;type:VARCHAR(50)"`                      // 操作名称
	BusinessType int64             `gorm:"column:business_type;type:BIGINT;NOT NULL"`              // 业务类型（1新增 2修改 3删除 4查询 5其它）
	Uri          string            `gorm:"column:uri;type:VARCHAR(100)"`                           // 请求地址
	OperIpAddr   string            `gorm:"column:oper_ip_addr;type:VARCHAR(50)"`                   // 主机地址
	OperLocation string            `gorm:"column:oper_location;type:VARCHAR(255)"`                 // 操作地点
	Req          sql.NullString    `gorm:"column:req;type:TEXT"`                                   // 请求参数
	Resp         sql.NullString    `gorm:"column:resp;type:TEXT"`                                  // 返回参数
	Code         int64             `gorm:"column:code;type:BIGINT;default:200;NOT NULL"`           // 返回状态（200成功 其它失败）
	Msg          string            `gorm:"column:msg;type:VARCHAR(255)"`                           // 提示消息
	CreatedTime  time.Time         `gorm:"column:created_time;default:CURRENT_TIMESTAMP;NOT NULL"` // 操作时间
}

func (m *SysTenantOperLog) TableName() string {
	return "sys_tenant_oper_log"
}

// 租户下的邮箱配置
type SysTenantConfig struct {
	ID             int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode     stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` // 租户编码
	Email          *SysTenantEmail   `gorm:"embedded;embeddedPrefix:email_"`
	RegisterRoleID int64             `gorm:"column:register_role_id;type:BIGINT;NOT NULL"` //注册分配的角色id
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:ri_mi"`
}

type SysTenantEmail struct {
	From     string `gorm:"column:from;type:VARCHAR(50);NOT NULL"`     // 发件人  你自己要发邮件的邮箱
	Host     string `gorm:"column:host;type:VARCHAR(50);NOT NULL"`     // 服务器地址 例如 smtp.qq.com  请前往QQ或者你要发邮件的邮箱查看其smtp协议
	Secret   string `gorm:"column:secret;type:VARCHAR(50);NOT NULL"`   // 密钥    用于登录的密钥 最好不要用邮箱密码 去邮箱smtp申请一个用于登录的密钥
	Nickname string `gorm:"column:nickname;type:VARCHAR(50);NOT NULL"` // 昵称    发件人昵称 通常为自己的邮箱
	Port     int64  `gorm:"column:port;type:int;default:465"`          // 端口     请前往QQ或者你要发邮件的邮箱查看其smtp协议 大多为 465
	IsSSL    int64  `gorm:"column:is_ssl;type:int;default:2"`          // 是否SSL   是否开启SSL
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;index"`
}

func (m *SysTenantConfig) TableName() string {
	return "sys_tenant_config"
}
