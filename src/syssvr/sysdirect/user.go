package sysdirect

import (
	client "github.com/i-Things/things/src/syssvr/client/user"
	server "github.com/i-Things/things/src/syssvr/internal/server/user"

	clientMenu "github.com/i-Things/things/src/syssvr/client/menu"
	clientRole "github.com/i-Things/things/src/syssvr/client/role"
	serverMenu "github.com/i-Things/things/src/syssvr/internal/server/menu"
	serverRole "github.com/i-Things/things/src/syssvr/internal/server/role"

	clientLog "github.com/i-Things/things/src/syssvr/client/log"
	serverLog "github.com/i-Things/things/src/syssvr/internal/server/log"

	clientCommon "github.com/i-Things/things/src/syssvr/client/common"
	serverCommon "github.com/i-Things/things/src/syssvr/internal/server/common"
)

func NewUser() client.User {
	userSvc := GetCtxSvc()
	return client.NewDirectUser(userSvc, server.NewUserServer(userSvc))
}

func NewRole() clientRole.Role {
	userSvc := GetCtxSvc()
	return clientRole.NewDirectRole(userSvc, serverRole.NewRoleServer(userSvc))
}

func NewMenu() clientMenu.Menu {
	userSvc := GetCtxSvc()
	return clientMenu.NewDirectMenu(userSvc, serverMenu.NewMenuServer(userSvc))
}

func NewCommon() clientCommon.Common {
	userSvc := GetCtxSvc()
	return clientCommon.NewDirectCommon(userSvc, serverCommon.NewCommonServer(userSvc))
}

func NewLog() clientLog.Log {
	userSvc := GetCtxSvc()
	return clientLog.NewDirectLog(userSvc, serverLog.NewLogServer(userSvc))
}

