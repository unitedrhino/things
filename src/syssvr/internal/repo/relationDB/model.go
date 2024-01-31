package relationDB

import (
	"database/sql"
	"github.com/i-Things/things/shared/stores"
	"time"
)

// 示例
type SysExample struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
}

// 应用信息
type SysAppInfo struct {
	ID      int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`        // id编号
	Code    string `gorm:"column:code;uniqueIndex:code;type:VARCHAR(100);NOT NULL"` // 应用编码
	Name    string `gorm:"column:name;uniqueIndex:name;type:VARCHAR(100);NOT NULL"` //应用名称
	Type    string `gorm:"column:type;type:VARCHAR(100);default:web;NOT NULL"`      //应用类型 web:web页面  app:应用  mini:小程序
	Desc    string `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`                  //应用描述
	BaseUrl string `gorm:"column:base_url;type:VARCHAR(100);NOT NULL"`              //应用首页
	LogoUrl string `gorm:"column:logo_url;type:VARCHAR(100);NOT NULL"`              //应用logo地址
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:code;uniqueIndex:name"`
}

func (m *SysAppInfo) TableName() string {
	return "sys_app_info"
}

// 应用默认绑定的模块
type SysAppModule struct {
	ID         int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	AppCode    string `gorm:"column:app_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"`    // 应用编码 这里只关联主应用,主应用授权,子应用也授权了
	ModuleCode string `gorm:"column:module_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"` // 模块编码
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:tc_ac"`
	Module      *SysModuleInfo     `gorm:"foreignKey:Code;references:ModuleCode"`
	App         *SysAppInfo        `gorm:"foreignKey:Code;references:AppCode"`
}

func (m *SysAppModule) TableName() string {
	return "sys_app_module"
}

// 模块管理表 模块是菜单和接口的集合体
type SysModuleInfo struct {
	ID         int64            `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`       // 编号
	Code       string           `gorm:"column:code;uniqueIndex:code;NOT NULL;type:VARCHAR(50)"` // 编码
	Type       int64            `gorm:"column:type;type:BIGINT;default:1;NOT NULL"`             // 类型   1:web页面  2:应用  3:小程序
	SubType    int64            `gorm:"column:sub_type;type:BIGINT;default:1;NOT NULL"`         // 类型   1：微应用   2：iframe内嵌 3: 原生菜单
	Order      int64            `gorm:"column:order_num;type:BIGINT;default:1;NOT NULL"`        // 左侧table排序序号
	Name       string           `gorm:"column:name;type:VARCHAR(50);NOT NULL"`                  // 菜单名称
	Path       string           `gorm:"column:path;type:VARCHAR(64);NOT NULL"`                  // 系统的path
	Url        string           `gorm:"column:url;type:VARCHAR(200);NOT NULL"`                  // 页面
	Icon       string           `gorm:"column:icon;type:VARCHAR(64);NOT NULL"`                  // 图标
	Body       string           `gorm:"column:body;type:VARCHAR(1024)"`                         // 菜单自定义数据
	HideInMenu int64            `gorm:"column:hide_in_menu;type:BIGINT;default:2;NOT NULL"`     // 是否隐藏菜单 1-是 2-否
	Desc       string           `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`                 // 备注
	Menus      []*SysModuleMenu `gorm:"foreignKey:ModuleCode;references:Code"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:code"`
}

func (m *SysModuleInfo) TableName() string {
	return "sys_module_info"
}

// 菜单管理表
type SysModuleMenu struct {
	ID         int64            `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`   // 编号
	ModuleCode string           `gorm:"column:module_code;type:VARCHAR(50);NOT NULL"`       // 模块编码
	ParentID   int64            `gorm:"column:parent_id;type:BIGINT;default:1;NOT NULL"`    // 父菜单ID，一级菜单为1
	Type       int64            `gorm:"column:type;type:BIGINT;default:1;NOT NULL"`         // 类型   1：菜单或者页面   2：iframe嵌入   3：外链跳转
	Order      int64            `gorm:"column:order_num;type:BIGINT;default:1;NOT NULL"`    // 左侧table排序序号
	Name       string           `gorm:"column:name;type:VARCHAR(50);NOT NULL"`              // 菜单名称
	Path       string           `gorm:"column:path;type:VARCHAR(64);NOT NULL"`              // 系统的path
	Component  string           `gorm:"column:component;type:VARCHAR(64);NOT NULL"`         // 页面
	Icon       string           `gorm:"column:icon;type:VARCHAR(64);NOT NULL"`              // 图标
	Redirect   string           `gorm:"column:redirect;type:VARCHAR(64);NOT NULL"`          // 路由重定向
	Body       string           `gorm:"column:body;type:VARCHAR(1024)"`                     // 菜单自定义数据
	HideInMenu int64            `gorm:"column:hide_in_menu;type:BIGINT;default:2;NOT NULL"` // 是否隐藏菜单 1-是 2-否
	Children   []*SysModuleMenu `gorm:"foreignKey:ID;references:ParentID"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;index"`
}

func (m *SysModuleMenu) TableName() string {
	return "sys_module_menu"
}

// 用户登录信息表
type SysUserInfo struct {
	TenantCode stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL;uniqueIndex:tc_un;uniqueIndex:tc_doi;uniqueIndex:tc_email;uniqueIndex:tc_phone;uniqueIndex:tc_wui;uniqueIndex:tc_woi"` // 租户编码

	UserID         int64          `gorm:"column:user_id;primary_key;AUTO_INCREMENT;type:BIGINT;NOT NULL"` // 用户id
	UserName       sql.NullString `gorm:"column:user_name;uniqueIndex:tc_un;type:VARCHAR(20)"`            // 登录用户名
	NickName       string         `gorm:"column:nick_name;type:VARCHAR(60);NOT NULL"`                     // 用户的昵称
	Password       string         `gorm:"column:password;type:CHAR(32);NOT NULL"`                         // 登录密码
	Email          sql.NullString `gorm:"column:email;uniqueIndex:tc_email;type:VARCHAR(255)"`            // 邮箱
	Phone          sql.NullString `gorm:"column:phone;uniqueIndex:tc_phone;type:VARCHAR(20)"`             // 手机号
	WechatUnionID  sql.NullString `gorm:"column:wechat_union_id;uniqueIndex:tc_wui;type:VARCHAR(20)"`     // 微信union id
	WechatOpenID   sql.NullString `gorm:"column:wechat_open_id;uniqueIndex:tc_woi;type:VARCHAR(20)"`      // 微信union id
	DingTalkUserID sql.NullString `gorm:"column:ding_talk_user_id;uniqueIndex:tc_doi;type:VARCHAR(20)"`
	LastIP         string         `gorm:"column:last_ip;type:VARCHAR(40);NOT NULL"`            // 最后登录ip
	RegIP          string         `gorm:"column:reg_ip;type:VARCHAR(40);NOT NULL"`             // 注册ip
	Sex            int64          `gorm:"column:sex;type:SMALLINT;default:3;NOT NULL"`         // 用户的性别，值为1时是男性，值为2时是女性，其他值为未知
	City           string         `gorm:"column:city;type:VARCHAR(50);NOT NULL"`               // 用户所在城市
	Country        string         `gorm:"column:country;type:VARCHAR(50);NOT NULL"`            // 用户所在国家
	Province       string         `gorm:"column:province;type:VARCHAR(50);NOT NULL"`           // 用户所在省份
	Language       string         `gorm:"column:language;type:VARCHAR(50);NOT NULL"`           // 用户的语言，简体中文为zh_CN
	HeadImg        string         `gorm:"column:head_img;type:VARCHAR(256);NOT NULL"`          // 用户头像
	Role           int64          `gorm:"column:role;type:BIGINT;NOT NULL"`                    // 用户默认角色（默认使用该角色）
	IsAllData      int64          `gorm:"column:is_all_data;type:SMALLINT;default:1;NOT NULL"` // 是否所有数据权限（1是，2否）
	Roles          []*SysUserRole `gorm:"foreignKey:UserID;references:UserID"`
	Tenant         *SysTenantInfo `gorm:"foreignKey:Code;references:TenantCode"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:tc_un;uniqueIndex:tc_doi;uniqueIndex:tc_email;uniqueIndex:tc_phone;uniqueIndex:tc_wui;uniqueIndex:tc_woi"`
}

func (m *SysUserInfo) TableName() string {
	return "sys_user_info"
}

// 应用菜单关联表
type SysUserRole struct {
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`      // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL;"`         // 租户编码
	UserID     int64             `gorm:"column:user_id;uniqueIndex:ri_mi;NOT NULL;type:BIGINT"` // 用户ID
	RoleID     int64             `gorm:"column:role_id;uniqueIndex:ri_mi;NOT NULL;type:BIGINT"` // 角色ID
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:ri_mi"`
}

func (m *SysUserRole) TableName() string {
	return "sys_user_role"
}

// 登录日志管理
type SysLoginLog struct {
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

func (m *SysLoginLog) TableName() string {
	return "sys_login_log"
}

// 操作日志管理
type SysOperLog struct {
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

func (m *SysOperLog) TableName() string {
	return "sys_oper_log"
}
