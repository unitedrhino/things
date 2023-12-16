package sysdirect

import (
	client "github.com/i-Things/things/src/syssvr/client/usermanage"
	server "github.com/i-Things/things/src/syssvr/internal/server/usermanage"

	clientRole "github.com/i-Things/things/src/syssvr/client/rolemanage"
	serverRole "github.com/i-Things/things/src/syssvr/internal/server/rolemanage"

	clientMenu "github.com/i-Things/things/src/syssvr/client/menumanage"
	serverMenu "github.com/i-Things/things/src/syssvr/internal/server/menumanage"

	clientLog "github.com/i-Things/things/src/syssvr/client/log"
	serverLog "github.com/i-Things/things/src/syssvr/internal/server/log"

	clientCommon "github.com/i-Things/things/src/syssvr/client/common"
	serverCommon "github.com/i-Things/things/src/syssvr/internal/server/common"

	clientApi "github.com/i-Things/things/src/syssvr/client/apimanage"
	serverApi "github.com/i-Things/things/src/syssvr/internal/server/apimanage"

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

func NewMenu(runSvr bool) clientMenu.MenuManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientMenu.NewDirectMenuManage(svcCtx, serverMenu.NewMenuManageServer(svcCtx))
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

func NewApi(runSvr bool) clientApi.ApiManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientApi.NewDirectApiManage(svcCtx, serverApi.NewApiManageServer(svcCtx))
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
