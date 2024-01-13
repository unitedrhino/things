package relationDB

import (
	"database/sql"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm/clause"
	"net/http"
)

func Migrate(c conf.Database) error {
	if c.IsInitTable == false {
		return nil
	}
	db := stores.GetCommonConn(nil)
	var needInitColumn bool
	if !db.Migrator().HasTable(&SysTenantUserInfo{}) {
		//需要初始化表
		needInitColumn = true
	}
	err := db.AutoMigrate(
		&SysTenantUserInfo{},
		&SysTenantRoleInfo{},
		&SysTenantRoleMenu{},
		&SysTenantRoleModule{},
		&SysModuleMenu{},
		&SysTenantLoginLog{},
		&SysTenantOperLog{},
		&SysModuleApi{},
		&SysTenantRoleApi{},
		&SysAreaInfo{},
		&SysProjectInfo{},
		&SysUserArea{},
		&SysUserProject{},
		&SysAppInfo{},
		&SysTenantRoleApp{},
		&SysTenantUserRole{},
		&SysTenantInfo{},
		&SysTenantApp{},
		&SysTenantConfig{},
		&SysModuleInfo{},
		&SysAppModule{},
		&SysTenantAppMenu{},
		&SysTenantAppApi{},
		&SysTenantAppModule{},
	)
	if err != nil {
		return err
	}
	//if true {
	//	db = db.Clauses(clause.OnConflict{DoNothing: true})
	//	if err := db.CreateInBatches(&MigrateModuleApi, 100).Error; err != nil {
	//		return err
	//	}
	//	for _, v := range MigrateModuleApi {
	//		data := SysTenantAppApi{
	//			TempLateID:   v.ID,
	//			TenantCode:   def.TenantCodeDefault,
	//			AppCode:      def.AppCore,
	//			SysModuleApi: v,
	//		}
	//		data.ID = 0
	//		MigrateTenantAppApi = append(MigrateTenantAppApi, data)
	//	}
	//	if err := db.CreateInBatches(&MigrateTenantAppApi, 100).Error; err != nil {
	//		return err
	//	}
	//}

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

	if err := db.CreateInBatches(&MigrateRoleMenu, 100).Error; err != nil {
		return err
	}

	//if err := db.CreateInBatches(&MigrateRoleApi, 100).Error; err != nil {
	//	return err
	//}
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
	{
		if err := db.CreateInBatches(&MigrateModuleApi, 100).Error; err != nil {
			return err
		}
		for _, v := range MigrateModuleApi {
			data := SysTenantAppApi{
				TempLateID:   v.ID,
				TenantCode:   def.TenantCodeDefault,
				AppCode:      def.AppCore,
				SysModuleApi: v,
			}
			data.ID = 0
			MigrateTenantAppApi = append(MigrateTenantAppApi, data)
		}
		if err := db.CreateInBatches(&MigrateTenantAppApi, 100).Error; err != nil {
			return err
		}
	}
	{
		if err := db.CreateInBatches(&MigrateModuleMenu, 100).Error; err != nil {
			return err
		}
		for _, v := range MigrateModuleMenu {
			data := SysTenantAppMenu{
				TenantCode:    def.TenantCodeDefault,
				SysModuleMenu: v,
			}
			data.ID = 0
			MigrateTenantAppMenu = append(MigrateTenantAppMenu, data)
		}
		if err := db.CreateInBatches(&MigrateModuleMenu, 100).Error; err != nil {
			return err
		}
	}
	{
		if err := db.CreateInBatches(&MigrateAppModule, 100).Error; err != nil {
			return err
		}
		for _, v := range MigrateAppModule {
			MigrateTenantAppModule = append(MigrateTenantAppModule, SysTenantAppModule{
				TenantCode:   def.TenantCodeDefault,
				SysAppModule: v,
			})
		}
		if err := db.CreateInBatches(&MigrateTenantAppModule, 100).Error; err != nil {
			return err
		}
	}

	return nil
}

func init() {
	for i := int64(1); i <= 100; i++ {
		MigrateRoleMenu = append(MigrateRoleMenu, SysTenantRoleMenu{
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
	MigrateProjectInfo = []SysProjectInfo{}
	MigrateModuleInfo  = []SysModuleInfo{
		{Name: "系统管理", Code: def.ModuleSystemManage},
		{Name: "租户管理", Code: def.ModuleTenantManage},
		{Name: "物联网", Code: def.ModuleThings},
		{Name: "大屏", Code: def.ModuleView},
	}
	MigrateAppModule = []SysAppModule{
		{
			AppCode:    def.AppCore,
			ModuleCode: def.ModuleThings,
		},
		{
			AppCode:    def.AppCore,
			ModuleCode: def.ModuleSystemManage,
		},
		{
			AppCode:    def.AppCore,
			ModuleCode: def.ModuleTenantManage,
		},
		{
			AppCode:    def.AppCore,
			ModuleCode: def.ModuleView,
		},
	}
	MigrateTenantAppModule = []SysTenantAppModule{}
	MigrateTenantAppApi    = []SysTenantAppApi{}
	MigrateTenantAppMenu   = []SysTenantAppMenu{}
	MigrateTenantConfig    = []SysTenantConfig{
		{TenantCode: def.TenantCodeDefault, RegisterRoleID: 2, Email: &SysTenantEmail{
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
	MigrateUserInfo   = []SysTenantUserInfo{
		{TenantCode: def.TenantCodeDefault, UserID: adminUserID, UserName: sql.NullString{String: "administrator", Valid: true}, Password: "4f0fded4a38abe7a3ea32f898bb82298", Role: 1, NickName: "iThings管理员", IsAllData: def.True},
	}
	MigrateUserRole = []SysTenantUserRole{
		{TenantCode: def.TenantCodeDefault, UserID: adminUserID, RoleID: 1},
	}
	MigrateRoleInfo = []SysTenantRoleInfo{
		{ID: 1, TenantCode: def.TenantCodeDefault, Name: "admin"},
		{ID: 2, TenantCode: def.TenantCodeDefault, Name: "client", Desc: "C端用户"}}
	MigrateRoleMenu []SysTenantRoleMenu
	MigrateRoleApp  = []SysTenantRoleApp{
		{RoleID: 1, TenantCode: def.TenantCodeDefault, AppCode: def.AppCore},
	}
	MigrateAppInfo = []SysAppInfo{
		{Code: def.AppCore, Name: "中台"},
	}

	MigrateModuleMenu = []SysModuleMenu{
		{ID: 2, ParentID: 1, Type: 1, Order: 2, Name: "设备管理", Path: "/deviceManagers", Component: "./deviceManagers/index.tsx", Icon: "icon_data_01", Redirect: "", HideInMenu: def.False},
		{ID: 6, ParentID: 2, Type: 1, Order: 1, Name: "产品", Path: "/deviceManagers/product/index", Component: "./deviceManagers/product/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 7, ParentID: 2, Type: 1, Order: 1, Name: "产品详情", Path: "/deviceManagers/product/detail/:id", Component: "./deviceManagers/product/detail/index", Icon: "icon_system", Redirect: "", HideInMenu: def.True},
		{ID: 8, ParentID: 2, Type: 1, Order: 2, Name: "设备", Path: "/deviceManagers/device/index", Component: "./deviceManagers/device/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 9, ParentID: 2, Type: 1, Order: 2, Name: "设备详情", Path: "/deviceManagers/device/detail/:id/:name/:type", Component: "./deviceManagers/device/detail/index", Icon: "icon_system", Redirect: "", HideInMenu: def.True},
		{ID: 23, ParentID: 2, Type: 1, Order: 3, Name: "分组", Path: "/deviceManagers/group/index", Component: "./deviceManagers/group/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 24, ParentID: 2, Type: 1, Order: 3, Name: "分组详情", Path: "/deviceManagers/group/detail/:id", Component: "./deviceManagers/group/detail/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.True},

		{ID: 4, ParentID: 1, Type: 1, Order: 4, Name: "运维监控", Path: "/operationsMonitorings", Component: "./operationsMonitorings/index.tsx", Icon: "icon_hvac", Redirect: "", HideInMenu: def.False},
		{ID: 13, ParentID: 4, Type: 1, Order: 1, Name: "固件升级", Path: "/operationsMonitorings/firmwareUpgrade/index", Component: "./operationsMonitorings/firmwareUpgrade/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 15, ParentID: 4, Type: 1, Order: 3, Name: "资源管理", Path: "/operationsMonitorings/resourceManagement/index", Component: "./operationsMonitorings/resourceManagement/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 16, ParentID: 4, Type: 1, Order: 4, Name: "远程配置", Path: "/operationsMonitorings/remoteConfiguration/index", Component: "./operationsMonitorings/remoteConfiguration/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 18, ParentID: 4, Type: 1, Order: 6, Name: "在线调试", Path: "/operationsMonitorings/onlineDebug/index", Component: "./operationsMonitorings/onlineDebug/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},

		{ID: 25, ParentID: 4, Type: 1, Order: 7, Name: "日志服务", Path: "/operationsMonitorings/logService/index", Component: "./operationsMonitorings/logService/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 35, ParentID: 1, Type: 1, Order: 1, Name: "首页", Path: "/home", Component: "./home/index.tsx", Icon: "icon_dosing", Redirect: "", HideInMenu: def.False},

		//{ID: 43, ParentID: 1, Type: 1, Order: 5, Name: "告警管理", Path: "/alarmManagers", Component: "./alarmManagers/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.False},
		//{ID: 44, ParentID: 43, Type: 1, Order: 1, Name: "告警配置", Path: "/alarmManagers/alarmConfiguration/index", Component: "./alarmManagers/alarmConfiguration/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.False},
		//{ID: 53, ParentID: 43, Type: 1, Order: 5, Name: "新增告警配置", Path: "/alarmManagers/alarmConfiguration/save", Component: "./alarmManagers/alarmConfiguration/addAlarmConfig/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.True},
		//{ID: 54, ParentID: 43, Type: 1, Order: 5, Name: "告警日志", Path: "/alarmManagers/alarmConfiguration/log/detail/:id/:level", Component: "./alarmManagers/alarmLog/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.True},
		//{ID: 45, ParentID: 43, Type: 1, Order: 5, Name: "告警记录", Path: "/alarmManagers/alarmConfiguration/log", Component: "./alarmManagers/alarmRecord/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.False},
		//{ID: 50, ParentID: 1, Type: 1, Order: 5, Name: "规则引擎", Path: "/ruleEngine", Component: "./ruleEngine/index.tsx", Icon: "icon_dosing", Redirect: "", HideInMenu: def.False},
		//{ID: 51, ParentID: 50, Type: 1, Order: 1, Name: "场景联动", Path: "/ruleEngine/scene/index", Component: "./ruleEngine/scene/index.tsx", Icon: "icon_device", Redirect: "", HideInMenu: def.False},

		{ID: 60, ParentID: 3, Type: 2, Order: 1, Name: "内嵌", Path: "/systemManagers/iframe", Component: "https://www.douyu.com", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 61, ParentID: 3, Type: 3, Order: 1, Name: "外链", Path: "/systemManagers/links", Component: "https://ant.design", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 70, ParentID: 3, Type: 1, Order: 1, Name: "任务管理", Path: "/systemManagers/timed", Component: "./systemManagers/timed/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 71, ParentID: 70, Type: 1, Order: 1, Name: "任务组", Path: "/systemManagers/timed/group", Component: "./systemManagers/timed/group/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 72, ParentID: 70, Type: 1, Order: 1, Name: "任务组详情", Path: "/systemManagers/timed/group/detail/:id", Component: "./systemManagers/timed/group/detail/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.True},
		{ID: 73, ParentID: 70, Type: 1, Order: 1, Name: "任务", Path: "/systemManagers/timed/task", Component: "./systemManagers/timed/task/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 74, ParentID: 70, Type: 1, Order: 1, Name: "任务详情", Path: "/systemManagers/timed/task/detail/:id", Component: "./systemManagers/timed/task/detail/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.True},
		{ID: 38, ParentID: 3, Type: 1, Order: 5, Name: "日志管理", Path: "/systemManagers/log", Component: "./systemManagers/log/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 39, ParentID: 38, Type: 1, Order: 1, Name: "操作日志", Path: "/systemManagers/log/operationLog/index", Component: "./systemManagers/log/operationLog/index.tsx", Icon: "icon_dosing", Redirect: "", HideInMenu: def.False},
		{ID: 41, ParentID: 38, Type: 1, Order: 2, Name: "登录日志", Path: "/systemManagers/log/loginLog/index", Component: "./systemManagers/log/loginLog/index", Icon: "icon_heat", Redirect: "", HideInMenu: def.False},
		{ID: 42, ParentID: 3, Type: 1, Order: 4, Name: "接口管理", Path: "/systemManagers/api/index", Component: "./systemManagers/api/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 10, ParentID: 3, Type: 1, Order: 1, Name: "用户管理", Path: "/systemManagers/user/index", Component: "./systemManagers/user/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 11, ParentID: 3, Type: 1, Order: 2, Name: "角色管理", Path: "/systemManagers/role/index", Component: "./systemManagers/role/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 12, ParentID: 3, Type: 1, Order: 3, Name: "菜单列表", Path: "/systemManagers/menu/index", Component: "./systemManagers/menu/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 3, ParentID: 1, Type: 1, Order: 9, Name: "系统管理", Path: "/systemManagers", Component: "./systemManagers/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},

		//视频服务菜单项
		{ID: 63, ParentID: 1, Type: 1, Order: 2, Name: "视频服务", Path: "/videoManagers", Component: "./videoManagers", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
		{ID: 64, ParentID: 63, Type: 1, Order: 1, Name: "流服务管理", Path: "/videoManagers/vidsrvmgr/index", Component: "./videoManagers/vidsrvmgr/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
		{ID: 65, ParentID: 63, Type: 1, Order: 3, Name: "视频流广场", Path: "/videoManagers/plaza/index", Component: "./videoManagers/plaza/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
		{ID: 66, ParentID: 63, Type: 1, Order: 2, Name: "视频流管理", Path: "/videoManagers/vidstream/index", Component: "./videoManagers/vidstream/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
		{ID: 67, ParentID: 63, Type: 1, Order: 4, Name: "视频回放", Path: "/videoManagers/playback/index", Component: "./videoManagers/playback/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
		{ID: 68, ParentID: 63, Type: 1, Order: 2, Name: "录像计划", Path: "/videoManagers/recordplan/index", Component: "./videoManagers/recordplan/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
		{ID: 69, ParentID: 63, Type: 1, Order: 1, Name: "流服务详细", Path: "/videoManagers/vidsrvmgr/detail/:id", Component: "./videoManagers/vidsrvmgr/detail/index", Icon: "icon_heat", Redirect: "", HideInMenu: 1},
		{ID: 75, ParentID: 63, Type: 1, Order: 1, Name: "视频流详细", Path: "/videoManagers/vidstream/detail/:id", Component: "./videoManagers/vidstream/detail/index", Icon: "icon_heat", Redirect: "", HideInMenu: 1},
	}

	//
	//MigrateModuleMenu = []SysModuleMenu{
	//	{ID: 2, ModuleCode: def.AppCore, ParentID: 1, Type: 1, Order: 2, Name: "设备管理", Path: "/deviceManagers", Component: "./deviceManagers/index.tsx", Icon: "icon_data_01", Redirect: "", HideInMenu: def.False},
	//	{ID: 3, ModuleCode: def.AppCore, ParentID: 1, Type: 1, Order: 9, Name: "系统管理", Path: "/systemManagers", Component: "./systemManagers/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 4, ModuleCode: def.AppCore, ParentID: 1, Type: 1, Order: 4, Name: "运维监控", Path: "/operationsMonitorings", Component: "./operationsMonitorings/index.tsx", Icon: "icon_hvac", Redirect: "", HideInMenu: def.False},
	//	{ID: 6, ModuleCode: def.AppCore, ParentID: 2, Type: 1, Order: 1, Name: "产品", Path: "/deviceManagers/product/index", Component: "./deviceManagers/product/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 7, ModuleCode: def.AppCore, ParentID: 2, Type: 1, Order: 1, Name: "产品详情", Path: "/deviceManagers/product/detail/:id", Component: "./deviceManagers/product/detail/index", Icon: "icon_system", Redirect: "", HideInMenu: def.True},
	//	{ID: 8, ModuleCode: def.AppCore, ParentID: 2, Type: 1, Order: 2, Name: "设备", Path: "/deviceManagers/device/index", Component: "./deviceManagers/device/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 9, ModuleCode: def.AppCore, ParentID: 2, Type: 1, Order: 2, Name: "设备详情", Path: "/deviceManagers/device/detail/:id/:name/:type", Component: "./deviceManagers/device/detail/index", Icon: "icon_system", Redirect: "", HideInMenu: def.True},
	//	{ID: 10, ModuleCode: def.AppCore, ParentID: 3, Type: 1, Order: 1, Name: "用户管理", Path: "/systemManagers/user/index", Component: "./systemManagers/user/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 11, ModuleCode: def.AppCore, ParentID: 3, Type: 1, Order: 2, Name: "角色管理", Path: "/systemManagers/role/index", Component: "./systemManagers/role/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 12, ModuleCode: def.AppCore, ParentID: 3, Type: 1, Order: 3, Name: "菜单列表", Path: "/systemManagers/menu/index", Component: "./systemManagers/menu/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 13, ModuleCode: def.AppCore, ParentID: 4, Type: 1, Order: 1, Name: "固件升级", Path: "/operationsMonitorings/firmwareUpgrade/index", Component: "./operationsMonitorings/firmwareUpgrade/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 15, ModuleCode: def.AppCore, ParentID: 4, Type: 1, Order: 3, Name: "资源管理", Path: "/operationsMonitorings/resourceManagement/index", Component: "./operationsMonitorings/resourceManagement/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 16, ModuleCode: def.AppCore, ParentID: 4, Type: 1, Order: 4, Name: "远程配置", Path: "/operationsMonitorings/remoteConfiguration/index", Component: "./operationsMonitorings/remoteConfiguration/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 18, ModuleCode: def.AppCore, ParentID: 4, Type: 1, Order: 6, Name: "在线调试", Path: "/operationsMonitorings/onlineDebug/index", Component: "./operationsMonitorings/onlineDebug/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 23, ModuleCode: def.AppCore, ParentID: 2, Type: 1, Order: 3, Name: "分组", Path: "/deviceManagers/group/index", Component: "./deviceManagers/group/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 24, ModuleCode: def.AppCore, ParentID: 2, Type: 1, Order: 3, Name: "分组详情", Path: "/deviceManagers/group/detail/:id", Component: "./deviceManagers/group/detail/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.True},
	//	{ID: 25, ModuleCode: def.AppCore, ParentID: 4, Type: 1, Order: 7, Name: "日志服务", Path: "/operationsMonitorings/logService/index", Component: "./operationsMonitorings/logService/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 35, ModuleCode: def.AppCore, ParentID: 1, Type: 1, Order: 1, Name: "首页", Path: "/home", Component: "./home/index.tsx", Icon: "icon_dosing", Redirect: "", HideInMenu: def.False},
	//	{ID: 38, ModuleCode: def.AppCore, ParentID: 3, Type: 1, Order: 5, Name: "日志管理", Path: "/systemManagers/log", Component: "./systemManagers/log/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 39, ModuleCode: def.AppCore, ParentID: 38, Type: 1, Order: 1, Name: "操作日志", Path: "/systemManagers/log/operationLog/index", Component: "./systemManagers/log/operationLog/index.tsx", Icon: "icon_dosing", Redirect: "", HideInMenu: def.False},
	//	{ID: 41, ModuleCode: def.AppCore, ParentID: 38, Type: 1, Order: 2, Name: "登录日志", Path: "/systemManagers/log/loginLog/index", Component: "./systemManagers/log/loginLog/index", Icon: "icon_heat", Redirect: "", HideInMenu: def.False},
	//	{ID: 42, ModuleCode: def.AppCore, ParentID: 3, Type: 1, Order: 4, Name: "接口管理", Path: "/systemManagers/api/index", Component: "./systemManagers/api/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 43, ModuleCode: def.AppCore, ParentID: 1, Type: 1, Order: 5, Name: "告警管理", Path: "/alarmManagers", Component: "./alarmManagers/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.False},
	//	{ID: 44, ModuleCode: def.AppCore, ParentID: 43, Type: 1, Order: 1, Name: "告警配置", Path: "/alarmManagers/alarmConfiguration/index", Component: "./alarmManagers/alarmConfiguration/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.False},
	//	{ID: 53, ModuleCode: def.AppCore, ParentID: 43, Type: 1, Order: 5, Name: "新增告警配置", Path: "/alarmManagers/alarmConfiguration/save", Component: "./alarmManagers/alarmConfiguration/addAlarmConfig/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.True},
	//	{ID: 54, ModuleCode: def.AppCore, ParentID: 43, Type: 1, Order: 5, Name: "告警日志", Path: "/alarmManagers/alarmConfiguration/log/detail/:id/:level", Component: "./alarmManagers/alarmLog/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.True},
	//	{ID: 45, ModuleCode: def.AppCore, ParentID: 43, Type: 1, Order: 5, Name: "告警记录", Path: "/alarmManagers/alarmConfiguration/log", Component: "./alarmManagers/alarmRecord/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.False},
	//	{ID: 50, ModuleCode: def.AppCore, ParentID: 1, Type: 1, Order: 5, Name: "规则引擎", Path: "/ruleEngine", Component: "./ruleEngine/index.tsx", Icon: "icon_dosing", Redirect: "", HideInMenu: def.False},
	//	{ID: 51, ModuleCode: def.AppCore, ParentID: 50, Type: 1, Order: 1, Name: "场景联动", Path: "/ruleEngine/scene/index", Component: "./ruleEngine/scene/index.tsx", Icon: "icon_device", Redirect: "", HideInMenu: def.False},
	//	{ID: 60, ModuleCode: def.AppCore, ParentID: 3, Type: 2, Order: 1, Name: "内嵌", Path: "/systemManagers/iframe", Component: "https://www.douyu.com", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 61, ModuleCode: def.AppCore, ParentID: 3, Type: 3, Order: 1, Name: "外链", Path: "/systemManagers/links", Component: "https://ant.design", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 70, ModuleCode: def.AppCore, ParentID: 3, Type: 1, Order: 1, Name: "任务管理", Path: "/systemManagers/timed", Component: "./systemManagers/timed/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 71, ModuleCode: def.AppCore, ParentID: 70, Type: 1, Order: 1, Name: "任务组", Path: "/systemManagers/timed/group", Component: "./systemManagers/timed/group/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 72, ModuleCode: def.AppCore, ParentID: 70, Type: 1, Order: 1, Name: "任务组详情", Path: "/systemManagers/timed/group/detail/:id", Component: "./systemManagers/timed/group/detail/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.True},
	//	{ID: 73, ModuleCode: def.AppCore, ParentID: 70, Type: 1, Order: 1, Name: "任务", Path: "/systemManagers/timed/task", Component: "./systemManagers/timed/task/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 74, ModuleCode: def.AppCore, ParentID: 70, Type: 1, Order: 1, Name: "任务详情", Path: "/systemManagers/timed/task/detail/:id", Component: "./systemManagers/timed/task/detail/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.True},
	//	//视频服务菜单项
	//	{ID: 63, ModuleCode: def.AppCore, ParentID: 1, Type: 1, Order: 2, Name: "视频服务", Path: "/videoManagers", Component: "./videoManagers", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
	//	{ID: 64, ModuleCode: def.AppCore, ParentID: 63, Type: 1, Order: 1, Name: "流服务管理", Path: "/videoManagers/vidsrvmgr/index", Component: "./videoManagers/vidsrvmgr/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
	//	{ID: 65, ModuleCode: def.AppCore, ParentID: 63, Type: 1, Order: 3, Name: "视频流广场", Path: "/videoManagers/plaza/index", Component: "./videoManagers/plaza/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
	//	{ID: 66, ModuleCode: def.AppCore, ParentID: 63, Type: 1, Order: 2, Name: "视频流管理", Path: "/videoManagers/vidstream/index", Component: "./videoManagers/vidstream/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
	//	{ID: 67, ModuleCode: def.AppCore, ParentID: 63, Type: 1, Order: 4, Name: "视频回放", Path: "/videoManagers/playback/index", Component: "./videoManagers/playback/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
	//	{ID: 68, ModuleCode: def.AppCore, ParentID: 63, Type: 1, Order: 2, Name: "录像计划", Path: "/videoManagers/recordplan/index", Component: "./videoManagers/recordplan/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
	//	{ID: 69, ModuleCode: def.AppCore, ParentID: 63, Type: 1, Order: 1, Name: "流服务详细", Path: "/videoManagers/vidsrvmgr/detail/:id", Component: "./videoManagers/vidsrvmgr/detail/index", Icon: "icon_heat", Redirect: "", HideInMenu: 1},
	//	{ID: 75, ModuleCode: def.AppCore, ParentID: 63, Type: 1, Order: 1, Name: "视频流详细", Path: "/videoManagers/vidstream/detail/:id", Component: "./videoManagers/vidstream/detail/index", Icon: "icon_heat", Redirect: "", HideInMenu: 1},
	//}
	MigrateModuleApi = []SysModuleApi{
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/role/info/create", Method: http.MethodPost, Name: "添加角色", BusinessType: 1, Desc: ``, Group: "角色管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/schema/index", Method: http.MethodPost, Name: "获取产品物模型列表", BusinessType: 4, Desc: ``, Group: "物模型"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/auth/project/index", Method: http.MethodPost, Name: "获取用户项目权限列表", BusinessType: 4, Desc: ``, Group: "用户权限"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/stream/update", Method: http.MethodPost, Name: "更新流信息", BusinessType: 2, Desc: ``, Group: "视频流管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/remote-config/push-all", Method: http.MethodPost, Name: "推送配置", BusinessType: 5, Desc: ``, Group: "远程配置"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/info/create", Method: http.MethodPost, Name: "新增设备", BusinessType: 1, Desc: ``, Group: "设备管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/role/app/index", Method: http.MethodPost, Name: "获取角色对应应用列表", BusinessType: 4, Desc: ``, Group: "角色管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/app/menu/create", Method: http.MethodPost, Name: "添加菜单", BusinessType: 1, Desc: ``, Group: "菜单管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/info/multi-import", Method: http.MethodPost, Name: "批量导入设备", BusinessType: 5, Desc: `#### 前端处理逻辑建议：
- UI text 显示 导入成功 设备数：total - len(errdata)
- UI text 显示 导入失败 设备数：len(errdata)
- UI table 显示 导入失败设备清单明细`, Group: "设备管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/msg/event-log/index", Method: http.MethodPost, Name: "获取物模型事件历史记录", BusinessType: 4, Desc: ``, Group: "设备消息"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/msg/property-log/index", Method: http.MethodPost, Name: "获取单个id属性历史记录", BusinessType: 4, Desc: ``, Group: "设备消息"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/module/menu/index", Method: http.MethodPost, Name: "获取菜单列表", BusinessType: 4, Desc: ``, Group: "菜单"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/info/read", Method: http.MethodPost, Name: "获取设备详情", BusinessType: 4, Desc: ``, Group: "设备管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/job/cancel", Method: http.MethodPost, Name: "取消动态升级策略", BusinessType: 5, Desc: ``, Group: "升级批次管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/scene/info/delete", Method: http.MethodPost, Name: "删除场景信息", BusinessType: 3, Desc: ``, Group: "场景联动"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/firmware/device-info-read", Method: http.MethodPost, Name: "获取升级包可选设备信息,包含可用版本", BusinessType: 4, Desc: ``, Group: "升级包管理firmware"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/info/delete", Method: http.MethodPost, Name: "删除产品", BusinessType: 3, Desc: ``, Group: "产品管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/auth/project/multi-update", Method: http.MethodPost, Name: "授权用户项目权限（内部会先全删后重加）", BusinessType: 2, Desc: ``, Group: "用户权限"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/gbsip/createchn", Method: http.MethodPost, Name: "创建通道", BusinessType: 5, Desc: ``, Group: "国标协议服务"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/task/device-index", Method: http.MethodPost, Name: "批次设备列表", BusinessType: 4, Desc: ``, Group: "升级任务管理task"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/project/info/delete", Method: http.MethodPost, Name: "删除项目", BusinessType: 3, Desc: ``, Group: "项目管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/schema/create", Method: http.MethodPost, Name: "新增物模型功能", BusinessType: 1, Desc: ``, Group: "物模型"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/gateway/multi-create", Method: http.MethodPost, Name: "批量添加网关子设备", BusinessType: 1, Desc: ``, Group: "网关子设备管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/firmware/index", Method: http.MethodPost, Name: "获取升级包列表", BusinessType: 4, Desc: ``, Group: "升级包管理firmware"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/scene/info/manually-trigger", Method: http.MethodPost, Name: "手动触发场景联动", BusinessType: 5, Desc: ``, Group: "场景联动"},
		{ModuleCode: def.ModuleView, IsNeedAuth: 1, Route: "/api/v1/view/go-view/project/create", Method: http.MethodPost, Name: "创建项目", BusinessType: 1, Desc: ``, Group: "项目管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/custom/read", Method: http.MethodPost, Name: "获取产品自定义信息", BusinessType: 4, Desc: `物联网平台通过定义一种物的描述语言来描述物模型模块和功能，称为TSL（Thing Specification Language）`, Group: "自定义"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/log/login/index", Method: http.MethodPost, Name: "获取登录日志列表", BusinessType: 4, Desc: ``, Group: "日志管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/task/read", Method: http.MethodPost, Name: "升级任务信息", BusinessType: 4, Desc: ``, Group: "升级任务管理task"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/app/module/multi-update", Method: http.MethodPost, Name: "批量更新应用绑定的模块", BusinessType: 2, Desc: ``, Group: "应用管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/register/dev", Method: http.MethodPost, Name: "未命名接口(因为这个接口导致生成的接口会有前缀)", BusinessType: 5, Desc: ``, Group: "设备鉴权"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/forget-pwd", Method: http.MethodPost, Name: "用戶忘记密码", BusinessType: 5, Desc: `注册接口`, Group: "用户管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/timed/task/group/read", Method: http.MethodPost, Name: "获取任务组详情", BusinessType: 4, Desc: ``, Group: "任务"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/schema/tsl-import", Method: http.MethodPost, Name: "导入物模型tsl", BusinessType: 5, Desc: `物联网平台通过定义一种物的描述语言来描述物模型模块和功能，称为TSL（Thing Specification Language）`, Group: "物模型"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/role/module/index", Method: http.MethodPost, Name: "获取角色对应模块列表 ", BusinessType: 4, Desc: ``, Group: "角色管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/scene/info/create", Method: http.MethodPost, Name: "创建场景信息", BusinessType: 1, Desc: ``, Group: "场景联动"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/stream/create", Method: http.MethodPost, Name: "创建流（拉流）", BusinessType: 1, Desc: ``, Group: "视频流管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/app/menu/update", Method: http.MethodPost, Name: "更新菜单", BusinessType: 2, Desc: ``, Group: "菜单管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/app/info/update", Method: http.MethodPost, Name: "更新应用", BusinessType: 2, Desc: ``, Group: "应用管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/module/api/index", Method: http.MethodPost, Name: "获取接口列表", BusinessType: 4, Desc: ``, Group: "接口"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/area/info/read", Method: http.MethodPost, Name: "获取项目区域详情（不含子节点）", BusinessType: 4, Desc: ``, Group: "区域管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/otaFirmware/create", Method: http.MethodPost, Name: "添加升级包", BusinessType: 1, Desc: ``, Group: "升级包管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/task/cancel", Method: http.MethodPost, Name: "取消所有升级中的任务", BusinessType: 5, Desc: ``, Group: "升级任务管理task"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/info/update", Method: http.MethodPost, Name: "更新设备", BusinessType: 2, Desc: ``, Group: "设备管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/otaFirmware/index", Method: http.MethodPost, Name: "升级包列表", BusinessType: 4, Desc: ``, Group: "升级包管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/common/upload-url/create", Method: http.MethodPost, Name: "获取文件上传地址", BusinessType: 1, Desc: `接口返回signed-url ,前端获取到该url后，往该url put上传文件`, Group: "通用功能"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/role/api/index", Method: http.MethodPost, Name: "获取角色对应接口列表", BusinessType: 4, Desc: ``, Group: "角色管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/info/index", Method: http.MethodPost, Name: "获取租户列表", BusinessType: 4, Desc: ``, Group: "租户管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/area/info/tree", Method: http.MethodPost, Name: "获取项目区域树", BusinessType: 5, Desc: ``, Group: "区域管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/project/info/update", Method: http.MethodPost, Name: "更新项目", BusinessType: 2, Desc: ``, Group: "项目管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/job/deviceIndex", Method: http.MethodPost, Name: "获取设备所在的升级包升级批次列表", BusinessType: 5, Desc: ``, Group: "升级批次管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/timed/task/group/update", Method: http.MethodPost, Name: "更新任务组", BusinessType: 2, Desc: ``, Group: "任务"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/ctrl/getsvr", Method: http.MethodPost, Name: "获取流服务状态", BusinessType: 5, Desc: ``, Group: "流服务交互"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/info/delete", Method: http.MethodPost, Name: "删除流服务器", BusinessType: 3, Desc: ``, Group: "流服务器管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/gbsip/readdev", Method: http.MethodPost, Name: "获取设备详细", BusinessType: 5, Desc: ``, Group: "国标协议服务"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/area/info/delete", Method: http.MethodPost, Name: "删除项目区域", BusinessType: 3, Desc: ``, Group: "区域管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/auth/area/multi-update", Method: http.MethodPost, Name: "授权用户区域权限（内部会先全删后重加)", BusinessType: 2, Desc: ``, Group: "用户权限"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/schema/delete", Method: http.MethodPost, Name: "删除物模型功能", BusinessType: 3, Desc: ``, Group: "物模型"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/alarm/deal-record/index", Method: http.MethodPost, Name: "获取告警处理记录列表", BusinessType: 4, Desc: ``, Group: "处理记录"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/gbsip/deletedev", Method: http.MethodPost, Name: "删除设备", BusinessType: 5, Desc: ``, Group: "国标协议服务"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/log/oper/index", Method: http.MethodPost, Name: "获取操作日志列表", BusinessType: 4, Desc: ``, Group: "日志管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/group/device/multi-delete", Method: http.MethodPost, Name: "删除分组设备(支持批量)", BusinessType: 3, Desc: ``, Group: "设备分组"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/interact/send-property", Method: http.MethodPost, Name: "调用设备属性", BusinessType: 5, Desc: ``, Group: "设备交互"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/app/module/delete", Method: http.MethodPost, Name: "删除绑定模块", BusinessType: 3, Desc: ``, Group: "模块管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/task/create", Method: http.MethodPost, Name: "创建升级任务", BusinessType: 1, Desc: ``, Group: "升级任务管理task"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/task/index", Method: http.MethodPost, Name: "获取升级批次任务列表", BusinessType: 4, Desc: ``, Group: "升级任务管理task"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/alarm/info/create", Method: http.MethodPost, Name: "新增告警", BusinessType: 1, Desc: ``, Group: "告警管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/gbsip/updatedev", Method: http.MethodPost, Name: "更新设备信息", BusinessType: 5, Desc: ``, Group: "国标协议服务"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/timed/task/info/index", Method: http.MethodPost, Name: "获取任务列表", BusinessType: 4, Desc: ``, Group: "任务"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/info/create", Method: http.MethodPost, Name: "新增产品", BusinessType: 1, Desc: ``, Group: "产品管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/app/info/delete", Method: http.MethodPost, Name: "删除应用", BusinessType: 3, Desc: ``, Group: "应用管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/alarm/info/update", Method: http.MethodPost, Name: "更新告警", BusinessType: 2, Desc: ``, Group: "告警管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/group/info/create", Method: http.MethodPost, Name: "创建分组", BusinessType: 1, Desc: ``, Group: "设备分组"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/module/info/update", Method: http.MethodPost, Name: "更新模块", BusinessType: 2, Desc: ``, Group: "模块"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/role/app/multi-update", Method: http.MethodPost, Name: "更新角色对应应用列表", BusinessType: 2, Desc: ``, Group: "角色管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/self/read", Method: http.MethodPost, Name: "用户获取自己的用户信息", BusinessType: 4, Desc: ``, Group: "用户管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/app/create", Method: http.MethodPost, Name: "新增租户下的应用绑定", BusinessType: 1, Desc: ``, Group: "应用管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/register", Method: http.MethodPost, Name: "用户注册", BusinessType: 5, Desc: `注册接口`, Group: "用户管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/auth/root-check", Method: http.MethodPost, Name: "鉴定mqtt账号root权限", BusinessType: 5, Desc: ``, Group: "设备鉴权"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/role/info/index", Method: http.MethodPost, Name: "获取角色列表", BusinessType: 4, Desc: ``, Group: "角色管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/msg/property-latest/index", Method: http.MethodPost, Name: "获取最新属性", BusinessType: 4, Desc: ``, Group: "设备消息"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/app/module/index", Method: http.MethodPost, Name: "获取模块绑定列表", BusinessType: 4, Desc: ``, Group: "模块管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/gbsip/readinfo", Method: http.MethodPost, Name: "获取服务详细", BusinessType: 5, Desc: ``, Group: "国标协议服务"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/module/menu/create", Method: http.MethodPost, Name: "添加菜单", BusinessType: 1, Desc: ``, Group: "菜单"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/app/api/update", Method: http.MethodPost, Name: "更新接口", BusinessType: 2, Desc: ``, Group: "接口管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/app/index", Method: http.MethodPost, Name: "获取租户下绑定的应用列表", BusinessType: 4, Desc: ``, Group: "应用管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/category/create", Method: http.MethodPost, Name: "新增产品品类", BusinessType: 1, Desc: ``, Group: "产品品类"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/info/create", Method: http.MethodPost, Name: "新增租户", BusinessType: 1, Desc: ``, Group: "租户管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/module/info/create", Method: http.MethodPost, Name: "添加模块", BusinessType: 1, Desc: ``, Group: "模块"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/interact/get-property-reply", Method: http.MethodPost, Name: "请求设备获取设备最新属性", BusinessType: 5, Desc: ``, Group: "设备交互"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/ctrl/setsvr", Method: http.MethodPost, Name: "修改流服务状态", BusinessType: 5, Desc: ``, Group: "流服务交互"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/msg/hub-log/index", Method: http.MethodPost, Name: "获取云端诊断日志", BusinessType: 4, Desc: ``, Group: "设备消息"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/info/read", Method: http.MethodPost, Name: "获取租户详情", BusinessType: 4, Desc: ``, Group: "租户管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/alarm/info/read", Method: http.MethodPost, Name: "获取告警详情", BusinessType: 4, Desc: ``, Group: "告警管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/alarm/record/index", Method: http.MethodPost, Name: "获取告警记录列表", BusinessType: 4, Desc: ``, Group: "告警记录"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/gbsip/indexchn", Method: http.MethodPost, Name: "获取通道列表", BusinessType: 5, Desc: ``, Group: "国标协议服务"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/job/dynamicCreate", Method: http.MethodPost, Name: "创建动态升级批次", BusinessType: 5, Desc: ``, Group: "升级批次管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/otaFirmware/delete", Method: http.MethodPost, Name: "删除升级包", BusinessType: 3, Desc: ``, Group: "升级包管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/role/info/delete", Method: http.MethodPost, Name: "删除角色", BusinessType: 3, Desc: ``, Group: "角色管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/firmware/create", Method: http.MethodPost, Name: "创建升级包版本", BusinessType: 1, Desc: ``, Group: "升级包管理firmware"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/group/info/index", Method: http.MethodPost, Name: "获取分组列表", BusinessType: 4, Desc: ``, Group: "设备分组"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/auth/access", Method: http.MethodPost, Name: "设备操作认证", BusinessType: 5, Desc: ``, Group: "设备鉴权"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/job/firmwareIndex", Method: http.MethodPost, Name: "获取升级包下的升级任务批次列表", BusinessType: 5, Desc: ``, Group: "升级批次管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/common/config", Method: http.MethodPost, Name: "获取系统配置", BusinessType: 5, Desc: ``, Group: "通用功能"},
		{ModuleCode: def.ModuleView, IsNeedAuth: 1, Route: "/api/v1/view/go-view/project/detail/read", Method: http.MethodPost, Name: "获取项目详情信息", BusinessType: 4, Desc: ``, Group: "项目详情"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/task/device-cancel", Method: http.MethodPost, Name: "取消单个设备升级", BusinessType: 5, Desc: ``, Group: "升级任务管理task"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/info/index", Method: http.MethodPost, Name: "获取用户信息列表", BusinessType: 4, Desc: ``, Group: "用户管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/flow/info/index", Method: http.MethodPost, Name: "获取流列表", BusinessType: 4, Desc: ``, Group: "流"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/app/info/create", Method: http.MethodPost, Name: "添加应用", BusinessType: 1, Desc: ``, Group: "应用管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/app/module/multi-create", Method: http.MethodPost, Name: "批量添加模块绑定", BusinessType: 1, Desc: ``, Group: "模块管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/self/resource/app/index", Method: http.MethodPost, Name: "获取用户应用列表", BusinessType: 4, Desc: ``, Group: "用户资源"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/alarm/info/delete", Method: http.MethodPost, Name: "删除告警", BusinessType: 3, Desc: ``, Group: "告警管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/flow/info/update", Method: http.MethodPost, Name: "修改流", BusinessType: 2, Desc: ``, Group: "流"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/category/read", Method: http.MethodPost, Name: "获取产品品类详情", BusinessType: 4, Desc: ``, Group: "产品品类"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/info/delete", Method: http.MethodPost, Name: "删除设备", BusinessType: 3, Desc: ``, Group: "设备管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/info/delete", Method: http.MethodPost, Name: "删除租户", BusinessType: 3, Desc: ``, Group: "租户管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/self/resource/read", Method: http.MethodPost, Name: "获取用户资源", BusinessType: 4, Desc: ``, Group: "用户资源"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/info/active", Method: http.MethodPost, Name: "激活流服务器", BusinessType: 5, Desc: ``, Group: "流服务器管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/info/read", Method: http.MethodPost, Name: "获取流服详细", BusinessType: 4, Desc: `{
  "vidmgrID":"1113459"
}`, Group: "流服务器管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/role/menu/multi-update", Method: http.MethodPost, Name: "更新角色对应菜单列表", BusinessType: 2, Desc: ``, Group: "角色管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/stream/index", Method: http.MethodPost, Name: "获取流列表", BusinessType: 4, Desc: ``, Group: "视频流管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/timed/task/group/create", Method: http.MethodPost, Name: "创建任务组", BusinessType: 1, Desc: ``, Group: "任务"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/alarm/deal-record/create", Method: http.MethodPost, Name: "新增告警处理记录", BusinessType: 1, Desc: ``, Group: "处理记录"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/otaFirmware/update", Method: http.MethodPost, Name: "更新升级包", BusinessType: 2, Desc: ``, Group: "升级包管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/info/index", Method: http.MethodPost, Name: "获取产品列表", BusinessType: 4, Desc: ``, Group: "产品管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/module/info/delete", Method: http.MethodPost, Name: "删除模块", BusinessType: 3, Desc: ``, Group: "模块"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/job/staticCreate", Method: http.MethodPost, Name: "创建静态升级批次", BusinessType: 5, Desc: ``, Group: "升级批次管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/firmware/read", Method: http.MethodPost, Name: "获取升级包详情", BusinessType: 4, Desc: ``, Group: "升级包管理firmware"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/captcha", Method: http.MethodPost, Name: "获取验证码", BusinessType: 5, Desc: ``, Group: "用户管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/info/count", Method: http.MethodPost, Name: "获取设备在线数", BusinessType: 5, Desc: ``, Group: "流服务器管理"},
		{ModuleCode: def.ModuleView, IsNeedAuth: 1, Route: "/api/v1/view/go-view/project/delete", Method: http.MethodPost, Name: "删除项目", BusinessType: 3, Desc: ``, Group: "项目管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/app/info/index", Method: http.MethodPost, Name: "获取应用列表", BusinessType: 4, Desc: ``, Group: "应用管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/timed/task/send", Method: http.MethodPost, Name: "执行任务", BusinessType: 5, Desc: ``, Group: "任务"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/app/api/tree", Method: http.MethodPost, Name: "获取接口树", BusinessType: 5, Desc: ``, Group: "接口管理"},
		{ModuleCode: def.ModuleView, IsNeedAuth: 1, Route: "/api/v1/view/go-view/project/update", Method: http.MethodPost, Name: "更新项目", BusinessType: 2, Desc: ``, Group: "项目管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/role/info/update", Method: http.MethodPost, Name: "更新角色", BusinessType: 2, Desc: ``, Group: "角色管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/project/info/create", Method: http.MethodPost, Name: "新增项目", BusinessType: 1, Desc: ``, Group: "项目管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/interact/send-action", Method: http.MethodPost, Name: "调用设备行为", BusinessType: 5, Desc: ``, Group: "设备交互"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/job/read", Method: http.MethodPost, Name: "查询指定升级批次的详情", BusinessType: 4, Desc: ``, Group: "升级批次管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/group/info/update", Method: http.MethodPost, Name: "更新分组信息", BusinessType: 2, Desc: ``, Group: "设备分组"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/app/delete", Method: http.MethodPost, Name: "删除租户下绑定的应用", BusinessType: 3, Desc: ``, Group: "应用管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/app/module/create", Method: http.MethodPost, Name: "添加模块绑定", BusinessType: 1, Desc: ``, Group: "模块管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/stream/read", Method: http.MethodPost, Name: "查询流详细", BusinessType: 4, Desc: ``, Group: "视频流管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/gateway/multi-delete", Method: http.MethodPost, Name: "批量解绑网关子设备", BusinessType: 3, Desc: ``, Group: "网关子设备管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/info/read", Method: http.MethodPost, Name: "获取用户信息", BusinessType: 4, Desc: ``, Group: "用户管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/msg/sdk-log/index", Method: http.MethodPost, Name: "获取设备本地日志", BusinessType: 4, Desc: `获取设备主动上传的sdk日志`, Group: "设备消息"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/info/read", Method: http.MethodPost, Name: "获取产品详情", BusinessType: 4, Desc: ``, Group: "产品管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/schema/tsl-read", Method: http.MethodPost, Name: "获取产品物模型tsl", BusinessType: 4, Desc: `物联网平台通过定义一种物的描述语言来描述物模型模块和功能，称为TSL（Thing Specification Language）`, Group: "物模型"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/firmware/delete", Method: http.MethodPost, Name: "删除升级包", BusinessType: 3, Desc: ``, Group: "升级包管理firmware"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/module/api/update", Method: http.MethodPost, Name: "更新接口", BusinessType: 2, Desc: ``, Group: "接口"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/timed/task/info/update", Method: http.MethodPost, Name: "更新任务", BusinessType: 2, Desc: ``, Group: "任务"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/interact/action-read", Method: http.MethodPost, Name: "获取调用设备行为的结果", BusinessType: 4, Desc: ``, Group: "设备交互"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/area/info/create", Method: http.MethodPost, Name: "新增项目区域", BusinessType: 1, Desc: ``, Group: "区域管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/gbsip/indexdev", Method: http.MethodPost, Name: "获取设备列表", BusinessType: 5, Desc: ``, Group: "国标协议服务"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/info/update", Method: http.MethodPost, Name: "更新租户", BusinessType: 2, Desc: ``, Group: "租户管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/category/delete", Method: http.MethodPost, Name: "删除产品品类", BusinessType: 3, Desc: ``, Group: "产品品类"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/info/update", Method: http.MethodPost, Name: "更新产品", BusinessType: 2, Desc: ``, Group: "产品管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/role/index", Method: http.MethodPost, Name: "获取用户角色信息列表", BusinessType: 4, Desc: ``, Group: "用户管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/schema/update", Method: http.MethodPost, Name: "更新物模型功能", BusinessType: 2, Desc: ``, Group: "物模型"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/app/api/create", Method: http.MethodPost, Name: "添加接口", BusinessType: 1, Desc: ``, Group: "接口管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/app/module/index", Method: http.MethodPost, Name: "获取应用绑定的模块列表", BusinessType: 4, Desc: ``, Group: "应用管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/remote-config/lastest-read", Method: http.MethodPost, Name: "获取最新配置", BusinessType: 4, Desc: ``, Group: "远程配置"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/firmware/update", Method: http.MethodPost, Name: "更新升级包", BusinessType: 2, Desc: ``, Group: "升级包管理firmware"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/role/multi-update", Method: http.MethodPost, Name: "用户角色信息批量更新", BusinessType: 2, Desc: ``, Group: "用户管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/gbsip/updatechn", Method: http.MethodPost, Name: "更新通道信息", BusinessType: 5, Desc: ``, Group: "国标协议服务"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/self/update", Method: http.MethodPost, Name: "更新用户基本数据", BusinessType: 2, Desc: ``, Group: "用户资源"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/module/menu/update", Method: http.MethodPost, Name: "更新菜单", BusinessType: 2, Desc: ``, Group: "菜单"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/gbsip/readchn", Method: http.MethodPost, Name: "获取通道详细", BusinessType: 5, Desc: ``, Group: "国标协议服务"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/scene/info/update", Method: http.MethodPost, Name: "更新场景信息", BusinessType: 2, Desc: ``, Group: "场景联动"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/alarm/log/index", Method: http.MethodPost, Name: "获取告警流水日志记录列表", BusinessType: 4, Desc: ``, Group: "告警日志"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/app/api/delete", Method: http.MethodPost, Name: "删除接口", BusinessType: 3, Desc: ``, Group: "接口管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/info/count", Method: http.MethodPost, Name: "设备统计详情", BusinessType: 5, Desc: ``, Group: "设备管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/group/info/read", Method: http.MethodPost, Name: "获取分组详情信息", BusinessType: 4, Desc: ``, Group: "设备分组"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/app/menu/index", Method: http.MethodPost, Name: "获取菜单列表", BusinessType: 4, Desc: ``, Group: "菜单管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/remote-config/create", Method: http.MethodPost, Name: "创建配置", BusinessType: 1, Desc: ``, Group: "远程配置"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/timed/task/group/index", Method: http.MethodPost, Name: "获取任务组列表", BusinessType: 4, Desc: ``, Group: "任务"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/group/device/index", Method: http.MethodPost, Name: "获取分组设备列表", BusinessType: 4, Desc: ``, Group: "设备分组"},
		{ModuleCode: def.ModuleView, IsNeedAuth: 1, Route: "/api/v1/view/go-view/project/index", Method: http.MethodPost, Name: "获取项目列表", BusinessType: 4, Desc: ``, Group: "项目管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/common/upload-file", Method: http.MethodPost, Name: "文件直传接口", BusinessType: 5, Desc: ``, Group: "通用功能"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/job/verify", Method: http.MethodPost, Name: "验证升级包", BusinessType: 5, Desc: ``, Group: "升级批次管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/stream/delete", Method: http.MethodPost, Name: "删除流", BusinessType: 3, Desc: ``, Group: "视频流管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/otaFirmware/read", Method: http.MethodPost, Name: "查询升级包", BusinessType: 4, Desc: ``, Group: "升级包管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/info/index", Method: http.MethodPost, Name: "获取设备列表", BusinessType: 4, Desc: ``, Group: "设备管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/role/module/multi-update", Method: http.MethodPost, Name: "更新角色对应模块列表 ", BusinessType: 2, Desc: ``, Group: "角色管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/flow/info/create", Method: http.MethodPost, Name: "创建流", BusinessType: 1, Desc: ``, Group: "流"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/role/api/multi-update", Method: http.MethodPost, Name: "更新角色对应接口列表", BusinessType: 2, Desc: ``, Group: "角色管理"},
		{ModuleCode: def.ModuleView, IsNeedAuth: 1, Route: "/api/v1/view/go-view/project/detail/update", Method: http.MethodPost, Name: "更新项目详情信息", BusinessType: 2, Desc: ``, Group: "项目详情"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/remote-config/index", Method: http.MethodPost, Name: "获取配置列表", BusinessType: 4, Desc: ``, Group: "远程配置"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/project/info/index", Method: http.MethodPost, Name: "获取项目列表", BusinessType: 4, Desc: ``, Group: "项目管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/module/info/index", Method: http.MethodPost, Name: "获取模块列表", BusinessType: 4, Desc: ``, Group: "模块"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/flow/info/delete", Method: http.MethodPost, Name: "删除流", BusinessType: 3, Desc: ``, Group: "流"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/info/update", Method: http.MethodPost, Name: "更新流服务器", BusinessType: 2, Desc: ``, Group: "流服务器管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/info/create", Method: http.MethodPost, Name: "创建用户信息", BusinessType: 1, Desc: ``, Group: "用户管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/ctrl/restart", Method: http.MethodPost, Name: "重流服务", BusinessType: 5, Desc: ``, Group: "流服务交互"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/alarm/scene/delete", Method: http.MethodPost, Name: "删除告警和场景的关联", BusinessType: 3, Desc: ``, Group: "场景联动关联"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/msg/shadow/index", Method: http.MethodPost, Name: "获取设备影子列表", BusinessType: 4, Desc: ``, Group: "设备消息"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/scene/info/read", Method: http.MethodPost, Name: "获取场景信息", BusinessType: 4, Desc: ``, Group: "场景联动"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/custom/update", Method: http.MethodPost, Name: "更新自定义信息", BusinessType: 2, Desc: ``, Group: "自定义"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/alarm/info/index", Method: http.MethodPost, Name: "获取告警信息列表", BusinessType: 4, Desc: ``, Group: "告警管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/timed/task/info/create", Method: http.MethodPost, Name: "创建任务", BusinessType: 1, Desc: ``, Group: "任务"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/auth/register", Method: http.MethodPost, Name: "设备动态注册", BusinessType: 5, Desc: ``, Group: "设备鉴权"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/info/update", Method: http.MethodPost, Name: "更新用户信息", BusinessType: 2, Desc: ``, Group: "用户管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/auth/login", Method: http.MethodPost, Name: "设备登录认证", BusinessType: 5, Desc: ``, Group: "设备鉴权"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/timed/task/cancel", Method: http.MethodPost, Name: "取消任务", BusinessType: 5, Desc: ``, Group: "任务"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/self/resource/module/index", Method: http.MethodPost, Name: "获取用户模块列表", BusinessType: 4, Desc: ``, Group: "用户资源"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/group/info/delete", Method: http.MethodPost, Name: "删除分组", BusinessType: 3, Desc: ``, Group: "设备分组"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/scene/info/index", Method: http.MethodPost, Name: "获取场景列表", BusinessType: 4, Desc: ``, Group: "场景联动"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/role/menu/index", Method: http.MethodPost, Name: "获取角色对应菜单列表", BusinessType: 4, Desc: ``, Group: "角色管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/gbsip/stopchn", Method: http.MethodPost, Name: "通道暂停播放", BusinessType: 5, Desc: ``, Group: "国标协议服务"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/module/api/delete", Method: http.MethodPost, Name: "删除接口", BusinessType: 3, Desc: ``, Group: "接口"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/self/resource/api/index", Method: http.MethodPost, Name: "获取用户接口列表", BusinessType: 4, Desc: ``, Group: "用户资源"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/timed/task/info/delete", Method: http.MethodPost, Name: "删除任务", BusinessType: 3, Desc: ``, Group: "任务"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/stream/count", Method: http.MethodPost, Name: "统计在线的流", BusinessType: 5, Desc: ``, Group: "视频流管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/category/index", Method: http.MethodPost, Name: "获取产品品类列表", BusinessType: 4, Desc: ``, Group: "产品品类"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/product/category/update", Method: http.MethodPost, Name: "更新产品品类", BusinessType: 2, Desc: ``, Group: "产品品类"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/task/device-retry", Method: http.MethodPost, Name: "重试单个设备升级", BusinessType: 5, Desc: ``, Group: "升级任务管理task"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/rule/alarm/scene/multi-update", Method: http.MethodPost, Name: "更新告警和场景的关联", BusinessType: 2, Desc: ``, Group: "场景联动关联"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/gbsip/deletechn", Method: http.MethodPost, Name: "删除通道", BusinessType: 5, Desc: ``, Group: "国标协议服务"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/app/api/index", Method: http.MethodPost, Name: "获取接口列表", BusinessType: 4, Desc: ``, Group: "接口管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/ota/task/analysis", Method: http.MethodPost, Name: "升级状态统计", BusinessType: 5, Desc: ``, Group: "升级任务管理task"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/module/menu/delete", Method: http.MethodPost, Name: "删除菜单", BusinessType: 3, Desc: ``, Group: "菜单"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/interact/property-read", Method: http.MethodPost, Name: "获取调用设备属性的结果", BusinessType: 4, Desc: ``, Group: "设备交互"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/module/api/create", Method: http.MethodPost, Name: "添加接口", BusinessType: 1, Desc: ``, Group: "接口"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/area/info/index", Method: http.MethodPost, Name: "获取项目区域列表", BusinessType: 4, Desc: ``, Group: "区域管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/info/index", Method: http.MethodPost, Name: "获取流服务器列表", BusinessType: 4, Desc: ``, Group: "流服务器管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/interact/multi-send-property", Method: http.MethodPost, Name: "批量调用设备属性", BusinessType: 5, Desc: ``, Group: "设备交互"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/gbsip/createdev", Method: http.MethodPost, Name: "创建设备", BusinessType: 5, Desc: ``, Group: "国标协议服务"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/login", Method: http.MethodPost, Name: "登录", BusinessType: 5, Desc: ``, Group: "用户管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/area/info/update", Method: http.MethodPost, Name: "更新项目区域", BusinessType: 2, Desc: ``, Group: "区域管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/interact/send-msg", Method: http.MethodPost, Name: "发送消息给设备", BusinessType: 5, Desc: ``, Group: "设备交互"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/timed/task/group/delete", Method: http.MethodPost, Name: "删除任务组", BusinessType: 3, Desc: ``, Group: "任务"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/project/info/read", Method: http.MethodPost, Name: "获取项目信息", BusinessType: 4, Desc: ``, Group: "项目管理"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/gbsip/playchn", Method: http.MethodPost, Name: "通道播放", BusinessType: 5, Desc: ``, Group: "国标协议服务"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/auth/area/index", Method: http.MethodPost, Name: "获取用户区域权限列表", BusinessType: 4, Desc: ``, Group: "用户权限"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/device/gateway/index", Method: http.MethodPost, Name: "获取子设备列表", BusinessType: 4, Desc: ``, Group: "网关子设备管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/module/api/tree", Method: http.MethodPost, Name: "获取接口树", BusinessType: 5, Desc: ``, Group: "接口"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/self/resource/menu/index", Method: http.MethodPost, Name: "获取用户菜单列表", BusinessType: 4, Desc: ``, Group: "用户资源"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/tenant/app/menu/delete", Method: http.MethodPost, Name: "删除菜单", BusinessType: 3, Desc: ``, Group: "菜单管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/user/info/delete", Method: http.MethodPost, Name: "删除用户", BusinessType: 3, Desc: ``, Group: "用户管理"},
		{ModuleCode: def.ModuleTenantManage, IsNeedAuth: 1, Route: "/api/v1/system/timed/task/info/read", Method: http.MethodPost, Name: "获取任务详情", BusinessType: 4, Desc: ``, Group: "任务"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/group/device/multi-create", Method: http.MethodPost, Name: "批量更新分组设备", BusinessType: 1, Desc: `会先删除后新增`, Group: "设备分组"},
		{ModuleCode: def.ModuleThings, IsNeedAuth: 1, Route: "/api/v1/things/vidmgr/info/create", Method: http.MethodPost, Name: "新增流服务器", BusinessType: 1, Desc: ``, Group: "流服务器管理"},
	}
	//MigrateModuleApi = []SysModuleApi{
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/product/info/update", Method: "POST", Name: "更新产品", BusinessType: 2, Desc: "", Group: "产品管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/product/info/create", Method: "POST", Name: "新增产品", BusinessType: 1, Desc: "", Group: "产品管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/product/info/read", Method: "POST", Name: "获取产品详情", BusinessType: 4, Desc: "", Group: "产品管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/product/info/delete", Method: "POST", Name: "删除产品", BusinessType: 3, Desc: "", Group: "产品管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/product/info/index", Method: "POST", Name: "获取产品列表", BusinessType: 4, Desc: "", Group: "产品管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/product/custom/read", Method: "POST", Name: "获取产品自定义信息", BusinessType: 4, Desc: "", Group: "产品自定义信息"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/product/custom/update", Method: "POST", Name: "更新产品自定义信息", BusinessType: 2, Desc: "", Group: "产品自定义信息"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/product/schema/index", Method: "POST", Name: "获取产品物模型列表", BusinessType: 4, Desc: "", Group: "物模型"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/product/schema/tsl-import", Method: "POST", Name: "导入物模型tsl", BusinessType: 1, Desc: "", Group: "物模型"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/product/schema/tsl-read", Method: "POST", Name: "获取产品物模型tsl", BusinessType: 4, Desc: "", Group: "物模型"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/product/schema/create", Method: "POST", Name: "新增物模型功能", BusinessType: 1, Desc: "", Group: "物模型"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/product/schema/update", Method: "POST", Name: "更新物模型功能", BusinessType: 2, Desc: "", Group: "物模型"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/product/schema/delete", Method: "POST", Name: "删除物模型功能", BusinessType: 3, Desc: "", Group: "物模型"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/product/remote-config/create", Method: "POST", Name: "创建配置", BusinessType: 1, Desc: "", Group: "产品远程配置"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/product/remote-config/index", Method: "POST", Name: "获取配置列表", BusinessType: 4, Desc: "", Group: "产品远程配置"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/product/remote-config/push-all", Method: "POST", Name: "推送配置", BusinessType: 5, Desc: "", Group: "产品远程配置"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/product/remote-config/lastest-read", Method: "POST", Name: "获取最新配置", BusinessType: 4, Desc: "", Group: "产品远程配置"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/group/info/create", Method: "POST", Name: "创建分组", BusinessType: 1, Desc: "", Group: "设备分组"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/group/info/index", Method: "POST", Name: "获取分组列表", BusinessType: 4, Desc: "", Group: "设备分组"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/group/info/read", Method: "POST", Name: "获取分组详情信息", BusinessType: 4, Desc: "", Group: "设备分组"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/group/info/update", Method: "POST", Name: "更新分组信息", BusinessType: 2, Desc: "", Group: "设备分组"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/group/info/delete", Method: "POST", Name: "删除分组", BusinessType: 3, Desc: "", Group: "设备分组"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/group/device/index", Method: "POST", Name: "获取分组设备列表", BusinessType: 4, Desc: "", Group: "设备分组"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/group/device/multi-create", Method: "POST", Name: "添加分组设备(支持批量)", BusinessType: 1, Desc: "", Group: "设备分组"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/group/device/multi-delete", Method: "POST", Name: "删除分组设备(支持批量)", BusinessType: 3, Desc: "", Group: "设备分组"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/info/index", Method: "POST", Name: "获取设备列表", BusinessType: 4, Desc: "", Group: "设备管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/info/read", Method: "POST", Name: "获取设备详情", BusinessType: 4, Desc: "", Group: "设备管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/info/create", Method: "POST", Name: "新增设备", BusinessType: 1, Desc: "", Group: "设备管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/info/delete", Method: "POST", Name: "删除设备", BusinessType: 3, Desc: "", Group: "设备管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/info/update", Method: "POST", Name: "更新设备", BusinessType: 2, Desc: "", Group: "设备管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/info/count", Method: "POST", Name: "设备统计详情", BusinessType: 4, Desc: "", Group: "设备管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/info/multi-import", Method: "POST", Name: "批量导入设备", BusinessType: 1, Desc: "", Group: "设备管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/auth/login", Method: "POST", Name: "设备登录认证", BusinessType: 5, Desc: "", Group: "设备鉴权"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/auth/root-check", Method: "POST", Name: "鉴定mqtt账号root权限", BusinessType: 5, Desc: "", Group: "设备鉴权"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/auth/access", Method: "POST", Name: "设备操作认证", BusinessType: 5, Desc: "", Group: "设备鉴权"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/msg/property-log/index", Method: "POST", Name: "获取单个id属性历史记录", BusinessType: 4, Desc: "", Group: "设备消息"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/msg/sdk-log/index", Method: "POST", Name: "获取设备本地日志", BusinessType: 4, Desc: "", Group: "设备消息"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/msg/hub-log/index", Method: "POST", Name: "获取云端诊断日志", BusinessType: 4, Desc: "", Group: "设备消息"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/msg/property-latest/index", Method: "POST", Name: "获取最新属性", BusinessType: 4, Desc: "", Group: "设备消息"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/msg/event-log/index", Method: "POST", Name: "获取物模型事件历史记录", BusinessType: 4, Desc: "", Group: "设备消息"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/interact/send-action", Method: "POST", Name: "同步调用设备行为", BusinessType: 5, Desc: "", Group: "设备交互"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/interact/send-property", Method: "POST", Name: "同步调用设备属性", BusinessType: 5, Desc: "", Group: "设备交互"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/interact/multi-send-property", Method: "POST", Name: "批量调用设备属性", BusinessType: 5, Desc: "", Group: "设备交互"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/interact/get-property-reply", Method: "POST", Name: "请求设备获取设备最新属性", BusinessType: 4, Desc: "", Group: "设备交互"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/interact/send-msg", Method: "POST", Name: "发送消息给设备", BusinessType: 5, Desc: "", Group: "设备交互"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/gateway/multi-create", Method: "POST", Name: "批量添加网关子设备", BusinessType: 1, Desc: "", Group: "网关子设备管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/gateway/multi-delete", Method: "POST", Name: "批量解绑网关子设备", BusinessType: 3, Desc: "", Group: "网关子设备管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/device/gateway/index", Method: "POST", Name: "获取子设备列表", BusinessType: 4, Desc: "", Group: "网关子设备管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/log/login/index", Method: "POST", Name: "获取登录日志列表", BusinessType: 4, Desc: "", Group: "日志管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/log/oper/index", Method: "POST", Name: "获取操作日志列表", BusinessType: 4, Desc: "", Group: "日志管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/role/info/create", Method: "POST", Name: "添加角色", BusinessType: 1, Desc: "", Group: "角色管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/role/info/index", Method: "POST", Name: "获取角色列表", BusinessType: 4, Desc: "", Group: "角色管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/role/info/update", Method: "POST", Name: "更新角色", BusinessType: 2, Desc: "", Group: "角色管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/role/info/delete", Method: "POST", Name: "删除角色", BusinessType: 3, Desc: "", Group: "角色管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/role/menu/multi-update", Method: "POST", Name: "更新角色对应菜单列表", BusinessType: 2, Desc: "", Group: "角色管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/role/menu/index", Method: "POST", Name: "获取角色对应菜单", BusinessType: 2, Desc: "", Group: "角色管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/role/app/multi-update", Method: "POST", Name: "更新角色对应应用列表", BusinessType: 2, Desc: "", Group: "角色管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/role/app/index", Method: "POST", Name: "获取角色对应应用", BusinessType: 2, Desc: "", Group: "角色管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/role/api/multi-update", Method: "POST", Name: "更新角色对应接口", BusinessType: 2, Desc: "", Group: "角色管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/role/api/index", Method: "POST", Name: "获取角色对应接口", BusinessType: 2, Desc: "", Group: "角色管理"},
	//
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/menu/info/create", Method: "POST", Name: "添加菜单", BusinessType: 1, Desc: "", Group: "菜单管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/menu/info/index", Method: "POST", Name: "获取菜单列表", BusinessType: 4, Desc: "", Group: "菜单管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/menu/info/update", Method: "POST", Name: "更新菜单", BusinessType: 2, Desc: "", Group: "菜单管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/menu/info/delete", Method: "POST", Name: "删除菜单", BusinessType: 3, Desc: "", Group: "菜单管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/user/info/create", Method: "POST", Name: "创建用户信息", BusinessType: 1, Desc: "", Group: "用户管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/user/captcha", Method: "POST", Name: "获取验证码", BusinessType: 5, Desc: "", Group: "用户管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/user/login", Method: "POST", Name: "登录", BusinessType: 5, Desc: "", Group: "用户管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/user/info/delete", Method: "POST", Name: "删除用户", BusinessType: 3, Desc: "", Group: "用户管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/user/info/read", Method: "POST", Name: "获取用户信息", BusinessType: 4, Desc: "", Group: "用户管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/user/info/update", Method: "POST", Name: "更新用户基本数据", BusinessType: 2, Desc: "", Group: "用户管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/user/info/index", Method: "POST", Name: "获取用户信息列表", BusinessType: 4, Desc: "", Group: "用户管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/user/role/index", Method: "POST", Name: "获取用户角色列表", BusinessType: 4, Desc: "", Group: "用户管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/user/role/multi-update", Method: "POST", Name: "获取用户角色列表", BusinessType: 4, Desc: "", Group: "用户管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/user/resource/read", Method: "POST", Name: "获取用户资源", BusinessType: 4, Desc: "", Group: "用户管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/common/config", Method: "POST", Name: "获取系统配置", BusinessType: 4, Desc: "", Group: "系统配置"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/app/info/create", Method: "POST", Name: "添加应用", BusinessType: 1, Desc: "", Group: "应用管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/app/info/index", Method: "POST", Name: "获取应用列表", BusinessType: 4, Desc: "", Group: "应用管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/app/info/read", Method: "POST", Name: "获取应用详情", BusinessType: 4, Desc: "", Group: "应用管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/app/info/update", Method: "POST", Name: "更新应用", BusinessType: 2, Desc: "", Group: "应用管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/app/info/delete", Method: "POST", Name: "删除应用", BusinessType: 3, Desc: "", Group: "应用管理"},
	//
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/tenant/info/create", Method: "POST", Name: "添加租户", BusinessType: 1, Desc: "", Group: "租户管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/tenant/info/index", Method: "POST", Name: "获取租户列表", BusinessType: 4, Desc: "", Group: "租户管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/tenant/info/read", Method: "POST", Name: "获取租户详情", BusinessType: 4, Desc: "", Group: "租户管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/tenant/info/update", Method: "POST", Name: "更新租户", BusinessType: 2, Desc: "", Group: "租户管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/tenant/info/delete", Method: "POST", Name: "删除租户", BusinessType: 3, Desc: "", Group: "租户管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/tenant/app/multi-update", Method: "POST", Name: "批量更新租户应用", BusinessType: 2, Desc: "", Group: "租户管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/tenant/app/index", Method: "POST", Name: "获取租户应用列表", BusinessType: 2, Desc: "", Group: "租户管理"},
	//
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/api/info/create", Method: "POST", Name: "添加接口", BusinessType: 1, Desc: "", Group: "接口管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/api/info/index", Method: "POST", Name: "获取接口列表", BusinessType: 4, Desc: "", Group: "接口管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/api/info/update", Method: "POST", Name: "更新接口", BusinessType: 2, Desc: "", Group: "接口管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/api/info/delete", Method: "POST", Name: "删除接口", BusinessType: 3, Desc: "", Group: "接口管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/auth/api/index", Method: "POST", Name: "获取API权限列表", BusinessType: 4, Desc: "", Group: "权限管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/auth/api/multiUpdate", Method: "POST", Name: "更新API权限", BusinessType: 2, Desc: "", Group: "权限管理"},
	//
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/timed/task/info/create", Method: "POST", Name: "新增任务", BusinessType: 1, Desc: "", Group: "定时任务管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/timed/task/info/update", Method: "POST", Name: "更新任务", BusinessType: 2, Desc: "", Group: "定时任务管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/timed/task/info/delete", Method: "POST", Name: "删除任务", BusinessType: 3, Desc: "", Group: "定时任务管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/timed/task/info/index", Method: "POST", Name: "获取任务信息列表", BusinessType: 4, Desc: "", Group: "定时任务管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/timed/task/info/read", Method: "POST", Name: "获取任务详情", BusinessType: 4, Desc: "", Group: "定时任务管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/timed/task/send", Method: "POST", Name: "执行任务", BusinessType: 5, Desc: "", Group: "定时任务管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/timed/task/group/create", Method: "POST", Name: "新增任务分组", BusinessType: 1, Desc: "", Group: "定时任务管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/timed/task/group/update", Method: "POST", Name: "更新任务分组", BusinessType: 2, Desc: "", Group: "定时任务管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/timed/task/group/delete", Method: "POST", Name: "删除任务分组", BusinessType: 3, Desc: "", Group: "定时任务管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/timed/task/group/index", Method: "POST", Name: "获取任务分组信息列表", BusinessType: 4, Desc: "", Group: "定时任务管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/system/timed/task/group/read", Method: "POST", Name: "获取任务分组详情", BusinessType: 4, Desc: "", Group: "定时任务管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/scene/info/read", Method: "POST", Name: "获取场景信息", BusinessType: 4, Desc: "", Group: "场景联动"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/scene/info/index", Method: "POST", Name: "获取场景列表", BusinessType: 4, Desc: "", Group: "场景联动"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/scene/info/create", Method: "POST", Name: "创建场景信息", BusinessType: 1, Desc: "", Group: "场景联动"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/scene/info/update", Method: "POST", Name: "更新场景信息", BusinessType: 2, Desc: "", Group: "场景联动"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/scene/info/delete", Method: "POST", Name: "删除场景信息", BusinessType: 3, Desc: "", Group: "场景联动"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/flow/info/index", Method: "POST", Name: "获取流列表", BusinessType: 4, Desc: "", Group: "流"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/flow/info/create", Method: "POST", Name: "创建流", BusinessType: 1, Desc: "", Group: "流"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/flow/info/update", Method: "POST", Name: "修改流", BusinessType: 2, Desc: "", Group: "流"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/flow/info/delete", Method: "POST", Name: "删除流", BusinessType: 3, Desc: "", Group: "流"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/alarm/info/create", Method: "POST", Name: "新增告警", BusinessType: 1, Desc: "", Group: "告警管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/alarm/info/update", Method: "POST", Name: "更新告警", BusinessType: 2, Desc: "", Group: "告警管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/alarm/info/delete", Method: "POST", Name: "删除告警", BusinessType: 3, Desc: "", Group: "告警管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/alarm/info/index", Method: "POST", Name: "获取告警信息列表", BusinessType: 4, Desc: "", Group: "告警管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/alarm/info/read", Method: "POST", Name: "获取告警详情", BusinessType: 4, Desc: "", Group: "告警管理"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/alarm/scene/delete", Method: "POST", Name: "删除告警和场景的关联", BusinessType: 3, Desc: "", Group: "场景联动"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/alarm/log/index", Method: "POST", Name: "获取告警流水日志记录列表", BusinessType: 4, Desc: "", Group: "告警日志"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/alarm/record/index", Method: "POST", Name: "获取告警记录列表", BusinessType: 4, Desc: "", Group: "告警记录"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/alarm/deal-record/create", Method: "POST", Name: "新增告警处理记录", BusinessType: 1, Desc: "", Group: "处理记录"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/alarm/deal-record/index", Method: "POST", Name: "获取告警处理记录列表", BusinessType: 4, Desc: "", Group: "处理记录"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/rule/alarm/scene/multi-update", Method: "POST", Name: "更新告警和场景的关联", BusinessType: 2, Desc: "", Group: "场景联动"},
	//	//视频服务API接口
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/vidmgr/info/active", Method: "POST", Name: "流服务激活", BusinessType: 1, Desc: "", Group: "视频服务"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/vidmgr/info/count", Method: "POST", Name: "流服务器统计", BusinessType: 1, Desc: "", Group: "视频服务"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/vidmgr/info/create", Method: "POST", Name: "新增流服务器", BusinessType: 1, Desc: "", Group: "视频服务"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/vidmgr/info/delete", Method: "POST", Name: "删除流服务器", BusinessType: 1, Desc: "", Group: "视频服务"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/vidmgr/info/index", Method: "POST", Name: "获取流服务器列表", BusinessType: 1, Desc: "", Group: "视频服务"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/vidmgr/info/read", Method: "POST", Name: "获取流服详细", BusinessType: 1, Desc: "", Group: "视频服务"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/vidmgr/info/update", Method: "POST", Name: "更新流服务器", BusinessType: 1, Desc: "", Group: "视频服务"},
	//	//视频流API接口
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/vidmgr/stream/count", Method: "POST", Name: "视频流统计", BusinessType: 1, Desc: "", Group: "视频服务"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/vidmgr/stream/create", Method: "POST", Name: "新增视频流(拉流)", BusinessType: 1, Desc: "", Group: "视频服务"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/vidmgr/stream/delete", Method: "POST", Name: "删除视频流", BusinessType: 1, Desc: "", Group: "视频服务"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/vidmgr/stream/index", Method: "POST", Name: "获取视频流列表", BusinessType: 1, Desc: "", Group: "视频服务"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/vidmgr/stream/read", Method: "POST", Name: "获取视频流详细", BusinessType: 1, Desc: "", Group: "视频服务"},
	//	{ModuleCode: def.AppCore, Route: "/api/v1/things/vidmgr/stream/update", Method: "POST", Name: "更新视频流", BusinessType: 1, Desc: "", Group: "视频服务"},
	//}
	//MigrateRoleApi = []SysTenantRoleApi{
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/info/update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/info/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/info/read", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/info/delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/info/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/schema/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/schema/tsl-import", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/schema/tsl-read", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/schema/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/schema/update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/schema/delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/remote-config/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/remote-config/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/remote-config/push-all", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/remote-config/lastest-read", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/custom/read", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/product/custom/update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/group/info/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/group/info/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/group/info/read", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/group/info/update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/group/info/delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/group/device/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/group/device/multi-create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/group/device/multi-delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/info/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/info/read", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/info/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/info/delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/info/update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/info/count", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/auth/login", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/info/multi-import", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/auth/root-check", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/auth/access", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/msg/property-log/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/msg/sdk-log/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/msg/hub-log/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/msg/property-latest/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/msg/event-log/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/interact/send-action", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/interact/send-property", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/interact/multi-send-property", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/interact/send-msg", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/gateway/multi-create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/gateway/multi-delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/device/gateway/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/scene/info/read", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/scene/info/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/scene/info/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/scene/info/update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/scene/info/delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/flow/info/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/flow/info/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/flow/info/update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/flow/info/delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/info/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/info/update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/info/delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/info/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/info/read", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/scene/delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/log/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/record/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/deal-record/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/deal-record/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/rule/alarm/scene/multi-update", V4: "POST", V5: ""},
	//	//视频服务API接口
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/info/active", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/info/count", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/info/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/info/delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/info/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/info/read", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/info/update", V4: "POST", V5: ""},
	//	//视频流API接口
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/stream/count", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/stream/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/stream/delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/stream/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/stream/read", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/things/vidmgr/stream/update", V4: "POST", V5: ""},
	//
	//	//系统管理接口
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/info/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/info/update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/info/delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/info/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/info/read", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/send", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/group/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/group/update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/group/delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/group/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/timed/task/group/read", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/log/login/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/log/oper/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/info/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/info/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/info/update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/info/delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/menu/multi-update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/menu/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/app/multi-update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/app/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/api/multi-update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/role/api/index", V4: "POST", V5: ""},
	//
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/menu/info/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/menu/info/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/menu/info/update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/menu/info/delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/info/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/captcha", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/login", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/info/delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/info/read", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/info/update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/info/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/resource/read", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/role/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/user/role/multi-update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/common/config", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/app/info/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/app/info/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/app/info/update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/app/info/delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/api/info/read", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/api/info/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/api/info/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/api/info/update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/api/info/delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/tenant/info/create", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/tenant/info/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/tenant/info/read", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/tenant/info/update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/tenant/info/delete", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/tenant/app/multi-update", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/tenant/app/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/auth/api/index", V4: "POST", V5: ""},
	//	{PType: "p", V0: "1", V1: def.TenantCodeDefault, V2: def.AppCore, V3: "/api/v1/system/auth/api/multiUpdate", V4: "POST", V5: ""},
	//}
)
