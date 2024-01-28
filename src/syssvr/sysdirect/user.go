package sysdirect

import (
	client "github.com/i-Things/things/src/syssvr/client/usermanage"
	server "github.com/i-Things/things/src/syssvr/internal/server/usermanage"

	clientRole "github.com/i-Things/things/src/syssvr/client/rolemanage"
	serverRole "github.com/i-Things/things/src/syssvr/internal/server/rolemanage"

	clientAccess "github.com/i-Things/things/src/syssvr/client/accessmanage"
	serverAccess "github.com/i-Things/things/src/syssvr/internal/server/accessmanage"

	clientModule "github.com/i-Things/things/src/syssvr/client/modulemanage"
	serverModule "github.com/i-Things/things/src/syssvr/internal/server/modulemanage"

	clientLog "github.com/i-Things/things/src/syssvr/client/log"
	serverLog "github.com/i-Things/things/src/syssvr/internal/server/log"

	clientCommon "github.com/i-Things/things/src/syssvr/client/common"
	serverCommon "github.com/i-Things/things/src/syssvr/internal/server/common"

	clientApp "github.com/i-Things/things/src/syssvr/client/appmanage"
	serverApp "github.com/i-Things/things/src/syssvr/internal/server/appmanage"

	clientTenant "github.com/i-Things/things/src/syssvr/client/tenantmanage"
	serverTenant "github.com/i-Things/things/src/syssvr/internal/server/tenantmanage"

	clientProject "github.com/i-Things/things/src/syssvr/client/projectmanage"
	serverProject "github.com/i-Things/things/src/syssvr/internal/server/projectmanage"

	clientArea "github.com/i-Things/things/src/syssvr/client/areamanage"
	serverArea "github.com/i-Things/things/src/syssvr/internal/server/areamanage"
)

func NewUser(runSvr bool) client.UserManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return client.NewDirectUserManage(svcCtx, server.NewUserManageServer(svcCtx))
}

func NewRole(runSvr bool) clientRole.RoleManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientRole.NewDirectRoleManage(svcCtx, serverRole.NewRoleManageServer(svcCtx))
}
func NewAccess(runSvr bool) clientAccess.AccessManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientAccess.NewDirectAccessManage(svcCtx, serverAccess.NewAccessManageServer(svcCtx))
}

func NewModule(runSvr bool) clientModule.ModuleManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientModule.NewDirectModuleManage(svcCtx, serverModule.NewModuleManageServer(svcCtx))
}

func NewCommon(runSvr bool) clientCommon.Common {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientCommon.NewDirectCommon(svcCtx, serverCommon.NewCommonServer(svcCtx))
}

func NewLog(runSvr bool) clientLog.Log {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientLog.NewDirectLog(svcCtx, serverLog.NewLogServer(svcCtx))
}

func NewApp(runSvr bool) clientApp.AppManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientApp.NewDirectAppManage(svcCtx, serverApp.NewAppManageServer(svcCtx))
}

func NewTenantManage(runSvr bool) clientTenant.TenantManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientTenant.NewDirectTenantManage(svcCtx, serverTenant.NewTenantManageServer(svcCtx))
}

func NewProjectManage(runSvr bool) clientProject.ProjectManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientProject.NewDirectProjectManage(svcCtx, serverProject.NewProjectManageServer(svcCtx))
}
func NewAreaManage(runSvr bool) clientArea.AreaManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientArea.NewDirectAreaManage(svcCtx, serverArea.NewAreaManageServer(svcCtx))
}
