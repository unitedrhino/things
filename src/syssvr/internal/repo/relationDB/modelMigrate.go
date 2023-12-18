package relationDB

import (
	"database/sql"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm/clause"
)

func Migrate(c conf.Database) error {
	if c.IsInitTable == false {
		return nil
	}
	db := stores.GetCommonConn(nil)
	var needInitColumn bool
	if !db.Migrator().HasTable(&SysUserInfo{}) {
		//需要初始化表
		needInitColumn = true
	}
	err := db.AutoMigrate(
		&SysUserInfo{},
		&SysRoleInfo{},
		&SysRoleMenu{},
		&SysMenuInfo{},
		&SysLoginLog{},
		&SysOperLog{},
		&SysApiInfo{},
		&SysRoleApi{},
		&SysAreaInfo{},
		&SysProjectInfo{},
		&SysUserAuthArea{},
		&SysUserAuthProject{},
		&SysAppInfo{},
		&SysRoleApp{},
		&SysUserRole{},
		&SysTenantInfo{},
		&SysTenantApp{},
		&SysTenantConfig{},
	)
	if err != nil {
		return err
	}
	if needInitColumn {
		return migrateTableColumn()
	}
	return err
}
func migrateTableColumn() error {
	db := stores.GetCommonConn(nil).Clauses(clause.OnConflict{DoNothing: true})
	if err := db.CreateInBatches(&MigrateUserInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateRoleInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateMenuInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateRoleMenu, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateApiInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateRoleApi, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateUserRole, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateRoleApp, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateAppInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateTenantInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateTenantApp, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateTenantConfig, 100).Error; err != nil {
		return err
	}
	return nil
}

func init() {
	for i := int64(1); i <= 100; i++ {
		MigrateRoleMenu = append(MigrateRoleMenu, SysRoleMenu{
			TenantCode: def.TenantCodeDefault,
			RoleID:     1,
			AppCode:    def.AppCore,
			MenuID:     i,
		})
	}
}

const (
	adminUserID = 1740358057038188544
)

// 子应用管理员可以配置自己子应用的角色
var (
	MigrateTenantConfig = []SysTenantConfig{
		{TenantCode: def.TenantCodeDefault, RegisterRoleID: 1, Email: &SysTenantEmail{
			From:     "godlei6@qq.com",
			Host:     "smtp.qq.com",
			Secret:   "xxx",
			Nickname: "验证码机器人",
			Port:     465,
			IsSSL:    def.True},
		},
	}
	MigrateTenantInfo = []SysTenantInfo{{Code: def.TenantCodeDefault, Name: "默认租户", AdminUserID: adminUserID}}
	MigrateTenantApp  = []SysTenantApp{{TenantCode: def.TenantCodeDefault, AppCode: def.AppCore}}
	MigrateUserInfo   = []SysUserInfo{
		{TenantCode: def.TenantCodeDefault, UserID: adminUserID, UserName: sql.NullString{String: "administrator", Valid: true}, Password: "4f0fded4a38abe7a3ea32f898bb82298", Role: 1, NickName: "iThings管理员", IsAllData: def.True},
	}
	MigrateUserRole = []SysUserRole{
		{TenantCode: def.TenantCodeDefault, UserID: adminUserID, RoleID: 1},
	}
	MigrateRoleInfo = []SysRoleInfo{{ID: 1, TenantCode: def.TenantCodeDefault, Name: "admin"}}
	MigrateRoleMenu []SysRoleMenu
	MigrateRoleApp  = []SysRoleApp{
		{RoleID: 1, TenantCode: def.TenantCodeDefault, AppCode: def.AppCore},
	}
	MigrateAppInfo = []SysAppInfo{
		{Code: def.AppCore, Name: "中台"},
	}

	MigrateMenuInfo = []SysMenuInfo{
		{ID: 2, AppCode: def.AppCore, ParentID: 1, Type: 1, Order: 2, Name: "设备管理", Path: "/deviceManagers", Component: "./deviceManagers/index.tsx", Icon: "icon_data_01", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 3, AppCode: def.AppCore, ParentID: 1, Type: 1, Order: 9, Name: "系统管理", Path: "/systemManagers", Component: "./systemManagers/index.tsx", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 4, AppCode: def.AppCore, ParentID: 1, Type: 1, Order: 4, Name: "运维监控", Path: "/operationsMonitorings", Component: "./operationsMonitorings/index.tsx", Icon: "icon_hvac", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 6, AppCode: def.AppCore, ParentID: 2, Type: 1, Order: 1, Name: "产品", Path: "/deviceManagers/product/index", Component: "./deviceManagers/product/index", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 7, AppCode: def.AppCore, ParentID: 2, Type: 1, Order: 1, Name: "产品详情", Path: "/deviceManagers/product/detail/:id", Component: "./deviceManagers/product/detail/index", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.True},
		{ID: 8, AppCode: def.AppCore, ParentID: 2, Type: 1, Order: 2, Name: "设备", Path: "/deviceManagers/device/index", Component: "./deviceManagers/device/index", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 9, AppCode: def.AppCore, ParentID: 2, Type: 1, Order: 2, Name: "设备详情", Path: "/deviceManagers/device/detail/:id/:name/:type", Component: "./deviceManagers/device/detail/index", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.True},
		{ID: 10, AppCode: def.AppCore, ParentID: 3, Type: 1, Order: 1, Name: "用户管理", Path: "/systemManagers/user/index", Component: "./systemManagers/user/index", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 11, AppCode: def.AppCore, ParentID: 3, Type: 1, Order: 2, Name: "角色管理", Path: "/systemManagers/role/index", Component: "./systemManagers/role/index", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 12, AppCode: def.AppCore, ParentID: 3, Type: 1, Order: 3, Name: "菜单列表", Path: "/systemManagers/menu/index", Component: "./systemManagers/menu/index", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 13, AppCode: def.AppCore, ParentID: 4, Type: 1, Order: 1, Name: "固件升级", Path: "/operationsMonitorings/firmwareUpgrade/index", Component: "./operationsMonitorings/firmwareUpgrade/index.tsx", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 15, AppCode: def.AppCore, ParentID: 4, Type: 1, Order: 3, Name: "资源管理", Path: "/operationsMonitorings/resourceManagement/index", Component: "./operationsMonitorings/resourceManagement/index.tsx", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 16, AppCode: def.AppCore, ParentID: 4, Type: 1, Order: 4, Name: "远程配置", Path: "/operationsMonitorings/remoteConfiguration/index", Component: "./operationsMonitorings/remoteConfiguration/index.tsx", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 18, AppCode: def.AppCore, ParentID: 4, Type: 1, Order: 6, Name: "在线调试", Path: "/operationsMonitorings/onlineDebug/index", Component: "./operationsMonitorings/onlineDebug/index.tsx", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 23, AppCode: def.AppCore, ParentID: 2, Type: 1, Order: 3, Name: "分组", Path: "/deviceManagers/group/index", Component: "./deviceManagers/group/index.tsx", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 24, AppCode: def.AppCore, ParentID: 2, Type: 1, Order: 3, Name: "分组详情", Path: "/deviceManagers/group/detail/:id", Component: "./deviceManagers/group/detail/index.tsx", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.True},
		{ID: 25, AppCode: def.AppCore, ParentID: 4, Type: 1, Order: 7, Name: "日志服务", Path: "/operationsMonitorings/logService/index", Component: "./operationsMonitorings/logService/index.tsx", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 35, AppCode: def.AppCore, ParentID: 1, Type: 1, Order: 1, Name: "首页", Path: "/home", Component: "./home/index.tsx", Icon: "icon_dosing", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 38, AppCode: def.AppCore, ParentID: 3, Type: 1, Order: 5, Name: "日志管理", Path: "/systemManagers/log", Component: "./systemManagers/log/index.tsx", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 39, AppCode: def.AppCore, ParentID: 38, Type: 1, Order: 1, Name: "操作日志", Path: "/systemManagers/log/operationLog/index", Component: "./systemManagers/log/operationLog/index.tsx", Icon: "icon_dosing", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 41, AppCode: def.AppCore, ParentID: 38, Type: 1, Order: 2, Name: "登录日志", Path: "/systemManagers/log/loginLog/index", Component: "./systemManagers/log/loginLog/index", Icon: "icon_heat", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 42, AppCode: def.AppCore, ParentID: 3, Type: 1, Order: 4, Name: "接口管理", Path: "/systemManagers/api/index", Component: "./systemManagers/api/index", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 43, AppCode: def.AppCore, ParentID: 1, Type: 1, Order: 5, Name: "告警管理", Path: "/alarmManagers", Component: "./alarmManagers/index", Icon: "icon_ap", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 44, AppCode: def.AppCore, ParentID: 43, Type: 1, Order: 1, Name: "告警配置", Path: "/alarmManagers/alarmConfiguration/index", Component: "./alarmManagers/alarmConfiguration/index", Icon: "icon_ap", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 53, AppCode: def.AppCore, ParentID: 43, Type: 1, Order: 5, Name: "新增告警配置", Path: "/alarmManagers/alarmConfiguration/save", Component: "./alarmManagers/alarmConfiguration/addAlarmConfig/index", Icon: "icon_ap", Redirect: "", BackgroundUrl: "", HideInMenu: def.True},
		{ID: 54, AppCode: def.AppCore, ParentID: 43, Type: 1, Order: 5, Name: "告警日志", Path: "/alarmManagers/alarmConfiguration/log/detail/:id/:level", Component: "./alarmManagers/alarmLog/index", Icon: "icon_ap", Redirect: "", BackgroundUrl: "", HideInMenu: def.True},
		{ID: 45, AppCode: def.AppCore, ParentID: 43, Type: 1, Order: 5, Name: "告警记录", Path: "/alarmManagers/alarmConfiguration/log", Component: "./alarmManagers/alarmRecord/index", Icon: "icon_ap", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 50, AppCode: def.AppCore, ParentID: 1, Type: 1, Order: 5, Name: "规则引擎", Path: "/ruleEngine", Component: "./ruleEngine/index.tsx", Icon: "icon_dosing", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 51, AppCode: def.AppCore, ParentID: 50, Type: 1, Order: 1, Name: "场景联动", Path: "/ruleEngine/scene/index", Component: "./ruleEngine/scene/index.tsx", Icon: "icon_device", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 60, AppCode: def.AppCore, ParentID: 3, Type: 2, Order: 1, Name: "内嵌", Path: "/systemManagers/iframe", Component: "https://www.douyu.com", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 61, AppCode: def.AppCore, ParentID: 3, Type: 3, Order: 1, Name: "外链", Path: "/systemManagers/links", Component: "https://ant.design", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 70, AppCode: def.AppCore, ParentID: 3, Type: 1, Order: 1, Name: "任务管理", Path: "/systemManagers/timed", Component: "./systemManagers/timed/index.tsx", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 71, AppCode: def.AppCore, ParentID: 70, Type: 1, Order: 1, Name: "任务组", Path: "/systemManagers/timed/group", Component: "./systemManagers/timed/group/index.tsx", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 72, AppCode: def.AppCore, ParentID: 70, Type: 1, Order: 1, Name: "任务组详情", Path: "/systemManagers/timed/group/detail/:id", Component: "./systemManagers/timed/group/detail/index.tsx", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.True},
		{ID: 73, AppCode: def.AppCore, ParentID: 70, Type: 1, Order: 1, Name: "任务", Path: "/systemManagers/timed/task", Component: "./systemManagers/timed/task/index.tsx", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.False},
		{ID: 74, AppCode: def.AppCore, ParentID: 70, Type: 1, Order: 1, Name: "任务详情", Path: "/systemManagers/timed/task/detail/:id", Component: "./systemManagers/timed/task/detail/index.tsx", Icon: "icon_system", Redirect: "", BackgroundUrl: "", HideInMenu: def.True},
		//视频服务菜单项
		{ID: 63, AppCode: def.AppCore, ParentID: 1, Type: 1, Order: 2, Name: "视频服务", Path: "/videoManagers", Component: "./videoManagers", Icon: "icon_heat", Redirect: "", BackgroundUrl: "", HideInMenu: 2},
		{ID: 64, AppCode: def.AppCore, ParentID: 63, Type: 1, Order: 1, Name: "流服务管理", Path: "/videoManagers/vidsrvmgr/index", Component: "./videoManagers/vidsrvmgr/index.tsx", Icon: "icon_heat", Redirect: "", BackgroundUrl: "", HideInMenu: 2},
		{ID: 65, AppCode: def.AppCore, ParentID: 63, Type: 1, Order: 3, Name: "视频流广场", Path: "/videoManagers/plaza/index", Component: "./videoManagers/plaza/index.tsx", Icon: "icon_heat", Redirect: "", BackgroundUrl: "", HideInMenu: 2},
		{ID: 66, AppCode: def.AppCore, ParentID: 63, Type: 1, Order: 2, Name: "视频流管理", Path: "/videoManagers/vidstream/index", Component: "./videoManagers/vidstream/index.tsx", Icon: "icon_heat", Redirect: "", BackgroundUrl: "", HideInMenu: 2},
		{ID: 67, AppCode: def.AppCore, ParentID: 63, Type: 1, Order: 4, Name: "视频回放", Path: "/videoManagers/playback/index", Component: "./videoManagers/playback/index.tsx", Icon: "icon_heat", Redirect: "", BackgroundUrl: "", HideInMenu: 2},
		{ID: 68, AppCode: def.AppCore, ParentID: 63, Type: 1, Order: 2, Name: "录像计划", Path: "/videoManagers/recordplan/index", Component: "./videoManagers/recordplan/index.tsx", Icon: "icon_heat", Redirect: "", BackgroundUrl: "", HideInMenu: 2},
		{ID: 69, AppCode: def.AppCore, ParentID: 63, Type: 1, Order: 1, Name: "流服务详细", Path: "/videoManagers/vidsrvmgr/detail/:id", Component: "./videoManagers/vidsrvmgr/detail/index", Icon: "icon_heat", Redirect: "", BackgroundUrl: "", HideInMenu: 1},
		{ID: 75, AppCode: def.AppCore, ParentID: 63, Type: 1, Order: 1, Name: "视频流详细", Path: "/videoManagers/vidstream/detail/:id", Component: "./videoManagers/vidstream/detail/index", Icon: "icon_heat", Redirect: "", BackgroundUrl: "", HideInMenu: 1},
	}
	MigrateApiInfo = []SysApiInfo{
		{AppCode: def.AppCore, Route: "/api/v1/things/product/info/update", Method: "POST", Name: "更新产品", BusinessType: 2, Desc: "", Group: "产品管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/product/info/create", Method: "POST", Name: "新增产品", BusinessType: 1, Desc: "", Group: "产品管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/product/info/read", Method: "POST", Name: "获取产品详情", BusinessType: 4, Desc: "", Group: "产品管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/product/info/delete", Method: "POST", Name: "删除产品", BusinessType: 3, Desc: "", Group: "产品管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/product/info/index", Method: "POST", Name: "获取产品列表", BusinessType: 4, Desc: "", Group: "产品管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/product/custom/read", Method: "POST", Name: "获取产品自定义信息", BusinessType: 4, Desc: "", Group: "产品自定义信息"},
		{AppCode: def.AppCore, Route: "/api/v1/things/product/custom/update", Method: "POST", Name: "更新产品自定义信息", BusinessType: 2, Desc: "", Group: "产品自定义信息"},
		{AppCode: def.AppCore, Route: "/api/v1/things/product/schema/index", Method: "POST", Name: "获取产品物模型列表", BusinessType: 4, Desc: "", Group: "物模型"},
		{AppCode: def.AppCore, Route: "/api/v1/things/product/schema/tsl-import", Method: "POST", Name: "导入物模型tsl", BusinessType: 1, Desc: "", Group: "物模型"},
		{AppCode: def.AppCore, Route: "/api/v1/things/product/schema/tsl-read", Method: "POST", Name: "获取产品物模型tsl", BusinessType: 4, Desc: "", Group: "物模型"},
		{AppCode: def.AppCore, Route: "/api/v1/things/product/schema/create", Method: "POST", Name: "新增物模型功能", BusinessType: 1, Desc: "", Group: "物模型"},
		{AppCode: def.AppCore, Route: "/api/v1/things/product/schema/update", Method: "POST", Name: "更新物模型功能", BusinessType: 2, Desc: "", Group: "物模型"},
		{AppCode: def.AppCore, Route: "/api/v1/things/product/schema/delete", Method: "POST", Name: "删除物模型功能", BusinessType: 3, Desc: "", Group: "物模型"},
		{AppCode: def.AppCore, Route: "/api/v1/things/product/remote-config/create", Method: "POST", Name: "创建配置", BusinessType: 1, Desc: "", Group: "产品远程配置"},
		{AppCode: def.AppCore, Route: "/api/v1/things/product/remote-config/index", Method: "POST", Name: "获取配置列表", BusinessType: 4, Desc: "", Group: "产品远程配置"},
		{AppCode: def.AppCore, Route: "/api/v1/things/product/remote-config/push-all", Method: "POST", Name: "推送配置", BusinessType: 5, Desc: "", Group: "产品远程配置"},
		{AppCode: def.AppCore, Route: "/api/v1/things/product/remote-config/lastest-read", Method: "POST", Name: "获取最新配置", BusinessType: 4, Desc: "", Group: "产品远程配置"},
		{AppCode: def.AppCore, Route: "/api/v1/things/group/info/create", Method: "POST", Name: "创建分组", BusinessType: 1, Desc: "", Group: "设备分组"},
		{AppCode: def.AppCore, Route: "/api/v1/things/group/info/index", Method: "POST", Name: "获取分组列表", BusinessType: 4, Desc: "", Group: "设备分组"},
		{AppCode: def.AppCore, Route: "/api/v1/things/group/info/read", Method: "POST", Name: "获取分组详情信息", BusinessType: 4, Desc: "", Group: "设备分组"},
		{AppCode: def.AppCore, Route: "/api/v1/things/group/info/update", Method: "POST", Name: "更新分组信息", BusinessType: 2, Desc: "", Group: "设备分组"},
		{AppCode: def.AppCore, Route: "/api/v1/things/group/info/delete", Method: "POST", Name: "删除分组", BusinessType: 3, Desc: "", Group: "设备分组"},
		{AppCode: def.AppCore, Route: "/api/v1/things/group/device/index", Method: "POST", Name: "获取分组设备列表", BusinessType: 4, Desc: "", Group: "设备分组"},
		{AppCode: def.AppCore, Route: "/api/v1/things/group/device/multi-create", Method: "POST", Name: "添加分组设备(支持批量)", BusinessType: 1, Desc: "", Group: "设备分组"},
		{AppCode: def.AppCore, Route: "/api/v1/things/group/device/multi-delete", Method: "POST", Name: "删除分组设备(支持批量)", BusinessType: 3, Desc: "", Group: "设备分组"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/info/index", Method: "POST", Name: "获取设备列表", BusinessType: 4, Desc: "", Group: "设备管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/info/read", Method: "POST", Name: "获取设备详情", BusinessType: 4, Desc: "", Group: "设备管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/info/create", Method: "POST", Name: "新增设备", BusinessType: 1, Desc: "", Group: "设备管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/info/delete", Method: "POST", Name: "删除设备", BusinessType: 3, Desc: "", Group: "设备管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/info/update", Method: "POST", Name: "更新设备", BusinessType: 2, Desc: "", Group: "设备管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/info/count", Method: "POST", Name: "设备统计详情", BusinessType: 4, Desc: "", Group: "设备管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/info/multi-import", Method: "POST", Name: "批量导入设备", BusinessType: 1, Desc: "", Group: "设备管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/auth/login", Method: "POST", Name: "设备登录认证", BusinessType: 5, Desc: "", Group: "设备鉴权"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/auth/root-check", Method: "POST", Name: "鉴定mqtt账号root权限", BusinessType: 5, Desc: "", Group: "设备鉴权"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/auth/access", Method: "POST", Name: "设备操作认证", BusinessType: 5, Desc: "", Group: "设备鉴权"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/msg/property-log/index", Method: "POST", Name: "获取单个id属性历史记录", BusinessType: 4, Desc: "", Group: "设备消息"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/msg/sdk-log/index", Method: "POST", Name: "获取设备本地日志", BusinessType: 4, Desc: "", Group: "设备消息"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/msg/hub-log/index", Method: "POST", Name: "获取云端诊断日志", BusinessType: 4, Desc: "", Group: "设备消息"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/msg/property-latest/index", Method: "POST", Name: "获取最新属性", BusinessType: 4, Desc: "", Group: "设备消息"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/msg/event-log/index", Method: "POST", Name: "获取物模型事件历史记录", BusinessType: 4, Desc: "", Group: "设备消息"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/interact/send-action", Method: "POST", Name: "同步调用设备行为", BusinessType: 5, Desc: "", Group: "设备交互"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/interact/send-property", Method: "POST", Name: "同步调用设备属性", BusinessType: 5, Desc: "", Group: "设备交互"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/interact/multi-send-property", Method: "POST", Name: "批量调用设备属性", BusinessType: 5, Desc: "", Group: "设备交互"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/interact/get-property-reply", Method: "POST", Name: "请求设备获取设备最新属性", BusinessType: 4, Desc: "", Group: "设备交互"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/interact/send-msg", Method: "POST", Name: "发送消息给设备", BusinessType: 5, Desc: "", Group: "设备交互"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/gateway/multi-create", Method: "POST", Name: "批量添加网关子设备", BusinessType: 1, Desc: "", Group: "网关子设备管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/gateway/multi-delete", Method: "POST", Name: "批量解绑网关子设备", BusinessType: 3, Desc: "", Group: "网关子设备管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/device/gateway/index", Method: "POST", Name: "获取子设备列表", BusinessType: 4, Desc: "", Group: "网关子设备管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/log/login/index", Method: "POST", Name: "获取登录日志列表", BusinessType: 4, Desc: "", Group: "日志管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/log/oper/index", Method: "POST", Name: "获取操作日志列表", BusinessType: 4, Desc: "", Group: "日志管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/role/info/create", Method: "POST", Name: "添加角色", BusinessType: 1, Desc: "", Group: "角色管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/role/info/index", Method: "POST", Name: "获取角色列表", BusinessType: 4, Desc: "", Group: "角色管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/role/info/update", Method: "POST", Name: "更新角色", BusinessType: 2, Desc: "", Group: "角色管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/role/info/delete", Method: "POST", Name: "删除角色", BusinessType: 3, Desc: "", Group: "角色管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/role/menu/multi-update", Method: "POST", Name: "更新角色对应菜单列表", BusinessType: 2, Desc: "", Group: "角色管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/role/menu/index", Method: "POST", Name: "获取角色对应菜单", BusinessType: 2, Desc: "", Group: "角色管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/role/app/multi-update", Method: "POST", Name: "更新角色对应应用列表", BusinessType: 2, Desc: "", Group: "角色管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/role/app/index", Method: "POST", Name: "获取角色对应应用", BusinessType: 2, Desc: "", Group: "角色管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/role/api/multi-update", Method: "POST", Name: "更新角色对应接口", BusinessType: 2, Desc: "", Group: "角色管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/role/api/index", Method: "POST", Name: "获取角色对应接口", BusinessType: 2, Desc: "", Group: "角色管理"},

		{AppCode: def.AppCore, Route: "/api/v1/system/menu/info/create", Method: "POST", Name: "添加菜单", BusinessType: 1, Desc: "", Group: "菜单管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/menu/info/index", Method: "POST", Name: "获取菜单列表", BusinessType: 4, Desc: "", Group: "菜单管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/menu/info/update", Method: "POST", Name: "更新菜单", BusinessType: 2, Desc: "", Group: "菜单管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/menu/info/delete", Method: "POST", Name: "删除菜单", BusinessType: 3, Desc: "", Group: "菜单管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/user/info/create", Method: "POST", Name: "创建用户信息", BusinessType: 1, Desc: "", Group: "用户管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/user/captcha", Method: "POST", Name: "获取验证码", BusinessType: 5, Desc: "", Group: "用户管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/user/login", Method: "POST", Name: "登录", BusinessType: 5, Desc: "", Group: "用户管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/user/info/delete", Method: "POST", Name: "删除用户", BusinessType: 3, Desc: "", Group: "用户管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/user/info/read", Method: "POST", Name: "获取用户信息", BusinessType: 4, Desc: "", Group: "用户管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/user/info/update", Method: "POST", Name: "更新用户基本数据", BusinessType: 2, Desc: "", Group: "用户管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/user/info/index", Method: "POST", Name: "获取用户信息列表", BusinessType: 4, Desc: "", Group: "用户管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/user/role/index", Method: "POST", Name: "获取用户角色列表", BusinessType: 4, Desc: "", Group: "用户管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/user/role/multi-update", Method: "POST", Name: "获取用户角色列表", BusinessType: 4, Desc: "", Group: "用户管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/user/resource/read", Method: "POST", Name: "获取用户资源", BusinessType: 4, Desc: "", Group: "用户管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/common/config", Method: "POST", Name: "获取系统配置", BusinessType: 4, Desc: "", Group: "系统配置"},
		{AppCode: def.AppCore, Route: "/api/v1/system/app/info/create", Method: "POST", Name: "添加应用", BusinessType: 1, Desc: "", Group: "应用管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/app/info/index", Method: "POST", Name: "获取应用列表", BusinessType: 4, Desc: "", Group: "应用管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/app/info/read", Method: "POST", Name: "获取应用详情", BusinessType: 4, Desc: "", Group: "应用管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/app/info/update", Method: "POST", Name: "更新应用", BusinessType: 2, Desc: "", Group: "应用管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/app/info/delete", Method: "POST", Name: "删除应用", BusinessType: 3, Desc: "", Group: "应用管理"},

		{AppCode: def.AppCore, Route: "/api/v1/system/tenant/info/create", Method: "POST", Name: "添加租户", BusinessType: 1, Desc: "", Group: "租户管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/tenant/info/index", Method: "POST", Name: "获取租户列表", BusinessType: 4, Desc: "", Group: "租户管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/tenant/info/read", Method: "POST", Name: "获取租户详情", BusinessType: 4, Desc: "", Group: "租户管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/tenant/info/update", Method: "POST", Name: "更新租户", BusinessType: 2, Desc: "", Group: "租户管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/tenant/info/delete", Method: "POST", Name: "删除租户", BusinessType: 3, Desc: "", Group: "租户管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/tenant/app/multi-update", Method: "POST", Name: "批量更新租户应用", BusinessType: 2, Desc: "", Group: "租户管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/tenant/app/index", Method: "POST", Name: "获取租户应用列表", BusinessType: 2, Desc: "", Group: "租户管理"},

		{AppCode: def.AppCore, Route: "/api/v1/system/api/info/create", Method: "POST", Name: "添加接口", BusinessType: 1, Desc: "", Group: "接口管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/api/info/index", Method: "POST", Name: "获取接口列表", BusinessType: 4, Desc: "", Group: "接口管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/api/info/update", Method: "POST", Name: "更新接口", BusinessType: 2, Desc: "", Group: "接口管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/api/info/delete", Method: "POST", Name: "删除接口", BusinessType: 3, Desc: "", Group: "接口管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/auth/api/index", Method: "POST", Name: "获取API权限列表", BusinessType: 4, Desc: "", Group: "权限管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/auth/api/multiUpdate", Method: "POST", Name: "更新API权限", BusinessType: 2, Desc: "", Group: "权限管理"},

		{AppCode: def.AppCore, Route: "/api/v1/system/timed/task/info/create", Method: "POST", Name: "新增任务", BusinessType: 1, Desc: "", Group: "定时任务管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/timed/task/info/update", Method: "POST", Name: "更新任务", BusinessType: 2, Desc: "", Group: "定时任务管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/timed/task/info/delete", Method: "POST", Name: "删除任务", BusinessType: 3, Desc: "", Group: "定时任务管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/timed/task/info/index", Method: "POST", Name: "获取任务信息列表", BusinessType: 4, Desc: "", Group: "定时任务管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/timed/task/info/read", Method: "POST", Name: "获取任务详情", BusinessType: 4, Desc: "", Group: "定时任务管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/timed/task/send", Method: "POST", Name: "执行任务", BusinessType: 5, Desc: "", Group: "定时任务管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/timed/task/group/create", Method: "POST", Name: "新增任务分组", BusinessType: 1, Desc: "", Group: "定时任务管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/timed/task/group/update", Method: "POST", Name: "更新任务分组", BusinessType: 2, Desc: "", Group: "定时任务管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/timed/task/group/delete", Method: "POST", Name: "删除任务分组", BusinessType: 3, Desc: "", Group: "定时任务管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/timed/task/group/index", Method: "POST", Name: "获取任务分组信息列表", BusinessType: 4, Desc: "", Group: "定时任务管理"},
		{AppCode: def.AppCore, Route: "/api/v1/system/timed/task/group/read", Method: "POST", Name: "获取任务分组详情", BusinessType: 4, Desc: "", Group: "定时任务管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/scene/info/read", Method: "POST", Name: "获取场景信息", BusinessType: 4, Desc: "", Group: "场景联动"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/scene/info/index", Method: "POST", Name: "获取场景列表", BusinessType: 4, Desc: "", Group: "场景联动"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/scene/info/create", Method: "POST", Name: "创建场景信息", BusinessType: 1, Desc: "", Group: "场景联动"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/scene/info/update", Method: "POST", Name: "更新场景信息", BusinessType: 2, Desc: "", Group: "场景联动"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/scene/info/delete", Method: "POST", Name: "删除场景信息", BusinessType: 3, Desc: "", Group: "场景联动"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/flow/info/index", Method: "POST", Name: "获取流列表", BusinessType: 4, Desc: "", Group: "流"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/flow/info/create", Method: "POST", Name: "创建流", BusinessType: 1, Desc: "", Group: "流"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/flow/info/update", Method: "POST", Name: "修改流", BusinessType: 2, Desc: "", Group: "流"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/flow/info/delete", Method: "POST", Name: "删除流", BusinessType: 3, Desc: "", Group: "流"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/alarm/info/create", Method: "POST", Name: "新增告警", BusinessType: 1, Desc: "", Group: "告警管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/alarm/info/update", Method: "POST", Name: "更新告警", BusinessType: 2, Desc: "", Group: "告警管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/alarm/info/delete", Method: "POST", Name: "删除告警", BusinessType: 3, Desc: "", Group: "告警管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/alarm/info/index", Method: "POST", Name: "获取告警信息列表", BusinessType: 4, Desc: "", Group: "告警管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/alarm/info/read", Method: "POST", Name: "获取告警详情", BusinessType: 4, Desc: "", Group: "告警管理"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/alarm/scene/delete", Method: "POST", Name: "删除告警和场景的关联", BusinessType: 3, Desc: "", Group: "场景联动"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/alarm/log/index", Method: "POST", Name: "获取告警流水日志记录列表", BusinessType: 4, Desc: "", Group: "告警日志"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/alarm/record/index", Method: "POST", Name: "获取告警记录列表", BusinessType: 4, Desc: "", Group: "告警记录"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/alarm/deal-record/create", Method: "POST", Name: "新增告警处理记录", BusinessType: 1, Desc: "", Group: "处理记录"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/alarm/deal-record/index", Method: "POST", Name: "获取告警处理记录列表", BusinessType: 4, Desc: "", Group: "处理记录"},
		{AppCode: def.AppCore, Route: "/api/v1/things/rule/alarm/scene/multi-update", Method: "POST", Name: "更新告警和场景的关联", BusinessType: 2, Desc: "", Group: "场景联动"},
		//视频服务API接口
		{AppCode: def.AppCore, Route: "/api/v1/things/vidmgr/info/active", Method: "POST", Name: "流服务激活", BusinessType: 1, Desc: "", Group: "视频服务"},
		{AppCode: def.AppCore, Route: "/api/v1/things/vidmgr/info/count", Method: "POST", Name: "流服务器统计", BusinessType: 1, Desc: "", Group: "视频服务"},
		{AppCode: def.AppCore, Route: "/api/v1/things/vidmgr/info/create", Method: "POST", Name: "新增流服务器", BusinessType: 1, Desc: "", Group: "视频服务"},
		{AppCode: def.AppCore, Route: "/api/v1/things/vidmgr/info/delete", Method: "POST", Name: "删除流服务器", BusinessType: 1, Desc: "", Group: "视频服务"},
		{AppCode: def.AppCore, Route: "/api/v1/things/vidmgr/info/index", Method: "POST", Name: "获取流服务器列表", BusinessType: 1, Desc: "", Group: "视频服务"},
		{AppCode: def.AppCore, Route: "/api/v1/things/vidmgr/info/read", Method: "POST", Name: "获取流服详细", BusinessType: 1, Desc: "", Group: "视频服务"},
		{AppCode: def.AppCore, Route: "/api/v1/things/vidmgr/info/update", Method: "POST", Name: "更新流服务器", BusinessType: 1, Desc: "", Group: "视频服务"},
		//视频流API接口
		{AppCode: def.AppCore, Route: "/api/v1/things/vidmgr/stream/count", Method: "POST", Name: "视频流统计", BusinessType: 1, Desc: "", Group: "视频服务"},
		{AppCode: def.AppCore, Route: "/api/v1/things/vidmgr/stream/create", Method: "POST", Name: "新增视频流(拉流)", BusinessType: 1, Desc: "", Group: "视频服务"},
		{AppCode: def.AppCore, Route: "/api/v1/things/vidmgr/stream/delete", Method: "POST", Name: "删除视频流", BusinessType: 1, Desc: "", Group: "视频服务"},
		{AppCode: def.AppCore, Route: "/api/v1/things/vidmgr/stream/index", Method: "POST", Name: "获取视频流列表", BusinessType: 1, Desc: "", Group: "视频服务"},
		{AppCode: def.AppCore, Route: "/api/v1/things/vidmgr/stream/read", Method: "POST", Name: "获取视频流详细", BusinessType: 1, Desc: "", Group: "视频服务"},
		{AppCode: def.AppCore, Route: "/api/v1/things/vidmgr/stream/update", Method: "POST", Name: "更新视频流", BusinessType: 1, Desc: "", Group: "视频服务"},
	}
	MigrateRoleApi = []SysRoleApi{
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/info/update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/info/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/info/read", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/info/delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/info/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/schema/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/schema/tsl-import", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/schema/tsl-read", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/schema/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/schema/update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/schema/delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/remote-config/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/remote-config/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/remote-config/push-all", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/remote-config/lastest-read", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/custom/read", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/custom/update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/group/info/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/group/info/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/group/info/read", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/group/info/update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/group/info/delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/group/device/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/group/device/multi-create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/group/device/multi-delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/info/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/info/read", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/info/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/info/delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/info/update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/info/count", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/auth/login", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/info/multi-import", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/auth/root-check", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/auth/access", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/msg/property-log/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/msg/sdk-log/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/msg/hub-log/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/msg/property-latest/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/msg/event-log/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/interact/send-action", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/interact/send-property", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/interact/multi-send-property", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/interact/send-msg", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/gateway/multi-create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/gateway/multi-delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/gateway/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/scene/info/read", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/scene/info/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/scene/info/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/scene/info/update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/scene/info/delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/flow/info/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/flow/info/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/flow/info/update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/flow/info/delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/info/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/info/update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/info/delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/info/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/info/read", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/scene/delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/log/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/record/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/deal-record/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/deal-record/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/scene/multi-update", V4: "POST", V5: ""},
		//视频服务API接口
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/info/active", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/info/count", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/info/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/info/delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/info/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/info/read", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/info/update", V4: "POST", V5: ""},
		//视频流API接口
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/stream/count", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/stream/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/stream/delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/stream/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/stream/read", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/stream/update", V4: "POST", V5: ""},

		//系统管理接口
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/info/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/info/update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/info/delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/info/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/info/read", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/send", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/group/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/group/update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/group/delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/group/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/group/read", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/log/login/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/log/oper/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/info/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/info/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/info/update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/info/delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/menu/multi-update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/menu/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/app/multi-update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/app/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/api/multi-update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/api/index", V4: "POST", V5: ""},

		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/menu/info/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/menu/info/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/menu/info/update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/menu/info/delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/info/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/captcha", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/login", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/info/delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/info/read", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/info/update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/info/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/resource/read", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/role/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/role/multi-update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/common/config", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/app/info/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/app/info/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/app/info/update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/app/info/delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/api/info/read", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/api/info/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/api/info/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/api/info/update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/api/info/delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/tenant/info/create", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/tenant/info/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/tenant/info/read", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/tenant/info/update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/tenant/info/delete", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/tenant/app/multi-update", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/tenant/app/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/auth/api/index", V4: "POST", V5: ""},
		{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/auth/api/multiUpdate", V4: "POST", V5: ""},
	}
)
