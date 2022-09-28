package sysdirect

import (
	client "github.com/i-Things/things/src/syssvr/client/user"
	server "github.com/i-Things/things/src/syssvr/internal/server/user"

	clientMenu "github.com/i-Things/things/src/syssvr/client/menu"
	clientRole "github.com/i-Things/things/src/syssvr/client/role"
	serverMenu "github.com/i-Things/things/src/syssvr/internal/server/menu"
	serverRole "github.com/i-Things/things/src/syssvr/internal/server/role"
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
