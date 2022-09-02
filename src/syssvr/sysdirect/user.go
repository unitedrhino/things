package sysdirect

import (
	client "github.com/i-Things/things/src/syssvr/client/user"
	server "github.com/i-Things/things/src/syssvr/internal/server/user"

	clientRole "github.com/i-Things/things/src/syssvr/client/role"
	serverRole "github.com/i-Things/things/src/syssvr/internal/server/role"
)

func NewUser(config *Config) client.User {
	userSvc := getCtxSvc(config)
	return client.NewDirectUser(userSvc, server.NewUserServer(userSvc))
}

func NewRole(config *Config) clientRole.Role {
	userSvc := getCtxSvc(config)
	return clientRole.NewDirectRole(userSvc, serverRole.NewRoleServer(userSvc))
}
