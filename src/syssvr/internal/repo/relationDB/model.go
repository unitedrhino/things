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

// 用户登录信息表
type SysUserInfo struct {
	UserID     int64          `gorm:"column:user_id;primary_key;AUTO_INCREMENT;type:BIGINT;NOT NULL"` // 用户id
	UserName   sql.NullString `gorm:"column:user_name;uniqueIndex;type:VARCHAR(20)"`                  // 登录用户名
	Password   string         `gorm:"column:password;type:CHAR(32);NOT NULL"`                         // 登录密码
	Email      sql.NullString `gorm:"column:email;uniqueIndex;type:VARCHAR(255)"`                     // 邮箱
	Phone      sql.NullString `gorm:"column:phone;uniqueIndex;type:VARCHAR(20)"`                      // 手机号
	Wechat     sql.NullString `gorm:"column:wechat;uniqueIndex;type:VARCHAR(20)"`                     // 微信union id
	LastIP     string         `gorm:"column:last_ip;type:VARCHAR(40);NOT NULL"`                       // 最后登录ip
	RegIP      string         `gorm:"column:reg_ip;type:VARCHAR(40);NOT NULL"`                        // 注册ip
	NickName   string         `gorm:"column:nick_name;type:VARCHAR(60);NOT NULL"`                     // 用户的昵称
	Sex        int64          `gorm:"column:sex;type:SMALLINT;default:3;NOT NULL"`                    // 用户的性别，值为1时是男性，值为2时是女性，其他值为未知
	City       string         `gorm:"column:city;type:VARCHAR(50);NOT NULL"`                          // 用户所在城市
	Country    string         `gorm:"column:country;type:VARCHAR(50);NOT NULL"`                       // 用户所在国家
	Province   string         `gorm:"column:province;type:VARCHAR(50);NOT NULL"`                      // 用户所在省份
	Language   string         `gorm:"column:language;type:VARCHAR(50);NOT NULL"`                      // 用户的语言，简体中文为zh_CN
	HeadImgUrl string         `gorm:"column:head_img_url;type:VARCHAR(256);NOT NULL"`                 // 用户头像
	Role       int64          `gorm:"column:role;type:BIGINT;NOT NULL"`                               // 用户角色
	IsAllData  int64          `gorm:"column:is_all_data;type:SMALLINT;default:1;NOT NULL"`            // 是否所有数据权限（1是，2否）
	stores.Time
}

func (m *SysUserInfo) TableName() string {
	return "sys_user_info"
}

// 角色管理表
type SysRoleInfo struct {
	ID     int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`   // id编号
	Name   string `gorm:"column:name;uniqueIndex;type:VARCHAR(100);NOT NULL"` // 角色名称
	Remark string `gorm:"column:remark;type:VARCHAR(100);NOT NULL"`           // 备注
	stores.Time
	Status int64          `gorm:"column:status;type:SMALLINT;default:1"` // 状态  1:启用,2:禁用
	Menus  []*SysRoleMenu `gorm:"foreignKey:role_id;references:id"`
}

func (m *SysRoleInfo) TableName() string {
	return "sys_role_info"
}

// 角色菜单关联表
type SysRoleMenu struct {
	ID     int64 `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`      // id编号
	RoleID int64 `gorm:"column:role_id;uniqueIndex:ri_mi;NOT NULL;type:BIGINT"` // 角色ID
	MenuID int64 `gorm:"column:menu_id;uniqueIndex:ri_mi;NOT NULL;type:BIGINT"` // 菜单ID
	stores.Time
}

func (m *SysRoleMenu) TableName() string {
	return "sys_role_menu"
}

// 菜单管理表
type SysMenuInfo struct {
	ID            int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`   // 编号
	ParentID      int64  `gorm:"column:parent_id;type:BIGINT;default:1;NOT NULL"`    // 父菜单ID，一级菜单为1
	Type          int64  `gorm:"column:type;type:BIGINT;default:1;NOT NULL"`         // 类型   1：目录   2：菜单   3：按钮
	Order         int64  `gorm:"column:order_num;type:BIGINT;default:1;NOT NULL"`    // 左侧table排序序号
	Name          string `gorm:"column:name;type:VARCHAR(50);NOT NULL"`              // 菜单名称
	Path          string `gorm:"column:path;type:VARCHAR(64);NOT NULL"`              // 系统的path
	Component     string `gorm:"column:component;type:VARCHAR(64);NOT NULL"`         // 页面
	Icon          string `gorm:"column:icon;type:VARCHAR(64);NOT NULL"`              // 图标
	Redirect      string `gorm:"column:redirect;type:VARCHAR(64);NOT NULL"`          // 路由重定向
	BackgroundUrl string `gorm:"column:background_url;type:VARCHAR(128);NOT NULL"`   // 后台地址
	HideInMenu    int64  `gorm:"column:hide_in_menu;type:BIGINT;default:2;NOT NULL"` // 是否隐藏菜单 1-是 2-否
	stores.Time
}

func (m *SysMenuInfo) TableName() string {
	return "sys_menu_info"
}

// 登录日志管理
type SysLoginLog struct {
	ID            int64     `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`       // 编号
	UserID        int64     `gorm:"column:user_id;type:BIGINT;NOT NULL"`                    // 用户id
	UserName      string    `gorm:"column:user_name;type:VARCHAR(50)"`                      // 登录账号
	IpAddr        string    `gorm:"column:ip_addr;type:VARCHAR(50)"`                        // 登录IP地址
	LoginLocation string    `gorm:"column:login_location;type:VARCHAR(100)"`                // 登录地点
	Browser       string    `gorm:"column:browser;type:VARCHAR(50)"`                        // 浏览器类型
	Os            string    `gorm:"column:os;type:VARCHAR(50)"`                             // 操作系统
	Code          int64     `gorm:"column:code;type:BIGINT;default:200;NOT NULL"`           // 登录状态（200成功 其它失败）
	Msg           string    `gorm:"column:msg;type:VARCHAR(255)"`                           // 提示消息
	CreatedTime   time.Time `gorm:"column:created_time;default:CURRENT_TIMESTAMP;NOT NULL"` // 登录时间
}

func (m *SysLoginLog) TableName() string {
	return "sys_login_log"
}

// 操作日志管理
type SysOperLog struct {
	ID           int64          `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`       // 编号
	OperUserID   int64          `gorm:"column:oper_user_id;type:BIGINT;NOT NULL"`               // 用户id
	OperUserName string         `gorm:"column:oper_user_name;type:VARCHAR(50)"`                 // 操作人员名称
	OperName     string         `gorm:"column:oper_name;type:VARCHAR(50)"`                      // 操作名称
	BusinessType int64          `gorm:"column:business_type;type:BIGINT;NOT NULL"`              // 业务类型（1新增 2修改 3删除 4查询 5其它）
	Uri          string         `gorm:"column:uri;type:VARCHAR(100)"`                           // 请求地址
	OperIpAddr   string         `gorm:"column:oper_ip_addr;type:VARCHAR(50)"`                   // 主机地址
	OperLocation string         `gorm:"column:oper_location;type:VARCHAR(255)"`                 // 操作地点
	Req          sql.NullString `gorm:"column:req;type:TEXT"`                                   // 请求参数
	Resp         sql.NullString `gorm:"column:resp;type:TEXT"`                                  // 返回参数
	Code         int64          `gorm:"column:code;type:BIGINT;default:200;NOT NULL"`           // 返回状态（200成功 其它失败）
	Msg          string         `gorm:"column:msg;type:VARCHAR(255)"`                           // 提示消息
	CreatedTime  time.Time      `gorm:"column:created_time;default:CURRENT_TIMESTAMP;NOT NULL"` // 操作时间
}

func (m *SysOperLog) TableName() string {
	return "sys_oper_log"
}

// 接口管理
type SysApiInfo struct {
	ID           int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`    // 编号
	Route        string `gorm:"column:route;uniqueIndex;type:VARCHAR(100);NOT NULL"` // 路由
	Method       int64  `gorm:"column:method;type:BIGINT;NOT NULL"`                  // 请求方式（1 GET 2 POST 3 HEAD 4 OPTIONS 5 PUT 6 DELETE 7 TRACE 8 CONNECT 9 其它）
	Name         string `gorm:"column:name;type:VARCHAR(100);NOT NULL"`              // 请求名称
	BusinessType int64  `gorm:"column:business_type;type:BIGINT;NOT NULL"`           // 业务类型（1新增 2修改 3删除 4查询 5其它）
	Group        string `gorm:"column:group;type:VARCHAR(100);NOT NULL"`             // 接口组
	Desc         string `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`              // 备注
	stores.Time
}

func (m *SysApiInfo) TableName() string {
	return "sys_api_info"
}

// api权限管理
type SysApiAuth struct {
	ID    int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"` // 编号
	PType string `gorm:"column:p_type;type:VARCHAR(255);NOT NULL"`         // 策略类型，即策略的分类，例如"p"表示主体（provider）访问资源（resource）的许可权，"g"表示主体（provider）之间的关系访问控制
	V0    string `gorm:"column:v0;type:VARCHAR(255);NOT NULL"`             // 策略中的第一个参数，通常用于表示资源的归属范围（即限制访问的对象），例如资源所属的机构、部门、业务线、地域等
	V1    string `gorm:"column:v1;type:VARCHAR(255);NOT NULL"`             // 策略中的第二个参数，通常用于表示主体（provider），即需要访问资源的用户或者服务
	V2    string `gorm:"column:v2;type:VARCHAR(255);NOT NULL"`             // 策略中的第三个参数，通常用于表示资源（resource），即需要进行访问的对象
	V3    string `gorm:"column:v3;type:VARCHAR(255);NOT NULL"`             // 策略中的第四个参数，通常用于表示访问操作（permission），例如 “read”, “write”, “execute” 等
	V4    string `gorm:"column:v4;type:VARCHAR(255);NOT NULL"`             // 策略中的第五个参数，通常用于表示资源的类型（object type），例如表示是文件或者数据库表等
	V5    string `gorm:"column:v5;type:VARCHAR(255);NOT NULL"`             // 策略中的第六个参数，通常用于表示扩展信息，例如 IP 地址、端口号等
}

func (m *SysApiAuth) TableName() string {
	return "sys_api_auth"
}
