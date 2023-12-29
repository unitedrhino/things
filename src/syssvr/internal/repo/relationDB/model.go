package relationDB

import "github.com/i-Things/things/shared/stores"

// 示例
type SysExample struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
}

// 应用信息
type SysAppInfo struct {
	ID      int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`        // id编号
	Code    string `gorm:"column:code;uniqueIndex:code;type:VARCHAR(100);NOT NULL"` // 应用编码
	Name    string `gorm:"column:name;uniqueIndex:name;type:VARCHAR(100);NOT NULL"` //应用名称
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
}

func (m *SysAppModule) TableName() string {
	return "sys_app_module"
}

// 模块管理表 模块是菜单和接口的集合体
type SysModuleInfo struct {
	ID         int64            `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`       // 编号
	Code       string           `gorm:"column:code;uniqueIndex:code;NOT NULL;type:VARCHAR(50)"` // 编码
	Type       int64            `gorm:"column:type;type:BIGINT;default:1;NOT NULL"`             // 类型   1：菜单或者页面   2：iframe嵌入   3：外链跳转
	Order      int64            `gorm:"column:order_num;type:BIGINT;default:1;NOT NULL"`        // 左侧table排序序号
	Name       string           `gorm:"column:name;type:VARCHAR(50);NOT NULL"`                  // 菜单名称
	Path       string           `gorm:"column:path;type:VARCHAR(64);NOT NULL"`                  // 系统的path
	Url        string           `gorm:"column:url;type:VARCHAR(200);NOT NULL"`                  // 页面
	Icon       string           `gorm:"column:icon;type:VARCHAR(64);NOT NULL"`                  // 图标
	Body       string           `gorm:"column:body;type:VARCHAR(1024)"`                         // 菜单自定义数据
	HideInMenu int64            `gorm:"column:hide_in_menu;type:BIGINT;default:2;NOT NULL"`     // 是否隐藏菜单 1-是 2-否
	Desc       string           `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`                 // 备注
	Apis       []*SysModuleApi  `gorm:"foreignKey:ModuleCode;references:Code"`
	Menus      []*SysModuleMenu `gorm:"foreignKey:ModuleCode;references:Code"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:code"`
}

func (m *SysModuleInfo) TableName() string {
	return "sys_module_info"
}

// 接口管理
type SysModuleApi struct {
	ID           int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                   // 编号
	ModuleCode   string `gorm:"column:module_code;uniqueIndex:app_route;type:VARCHAR(50);NOT NULL"` // 模块编码
	Route        string `gorm:"column:route;uniqueIndex:app_route;type:VARCHAR(100);NOT NULL"`      // 路由
	Method       string `gorm:"column:method;type:VARCHAR(50);NOT NULL"`                            // 请求方式（1 GET 2 POST 3 HEAD 4 OPTIONS 5 PUT 6 DELETE 7 TRACE 8 CONNECT 9 其它）
	Name         string `gorm:"column:name;type:VARCHAR(100);NOT NULL"`                             // 请求名称
	BusinessType int64  `gorm:"column:business_type;type:BIGINT;NOT NULL"`                          // 业务类型（1新增 2修改 3删除 4查询 5其它）
	Group        string `gorm:"column:group;type:VARCHAR(100);NOT NULL"`                            // 接口组
	IsNeedAuth   int64  `gorm:"column:is_need_auth;type:BIGINT;default:1;NOT NULL"`                 // 是否需要认证（1是 2否）
	Desc         string `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`                             // 备注
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:app_route"`
}

func (m *SysModuleApi) TableName() string {
	return "sys_module_api"
}

// 菜单管理表
type SysModuleMenu struct {
	ID         int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`   // 编号
	ModuleCode string `gorm:"column:module_code;type:VARCHAR(50);NOT NULL"`       // 模块编码
	ParentID   int64  `gorm:"column:parent_id;type:BIGINT;default:1;NOT NULL"`    // 父菜单ID，一级菜单为1
	Type       int64  `gorm:"column:type;type:BIGINT;default:1;NOT NULL"`         // 类型   1：菜单或者页面   2：iframe嵌入   3：外链跳转
	Order      int64  `gorm:"column:order_num;type:BIGINT;default:1;NOT NULL"`    // 左侧table排序序号
	Name       string `gorm:"column:name;type:VARCHAR(50);NOT NULL"`              // 菜单名称
	Path       string `gorm:"column:path;type:VARCHAR(64);NOT NULL"`              // 系统的path
	Component  string `gorm:"column:component;type:VARCHAR(64);NOT NULL"`         // 页面
	Icon       string `gorm:"column:icon;type:VARCHAR(64);NOT NULL"`              // 图标
	Redirect   string `gorm:"column:redirect;type:VARCHAR(64);NOT NULL"`          // 路由重定向
	Body       string `gorm:"column:body;type:VARCHAR(1024)"`                     // 菜单自定义数据
	HideInMenu int64  `gorm:"column:hide_in_menu;type:BIGINT;default:2;NOT NULL"` // 是否隐藏菜单 1-是 2-否
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;index"`
}

func (m *SysModuleMenu) TableName() string {
	return "sys_module_menu"
}
