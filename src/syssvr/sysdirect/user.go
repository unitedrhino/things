package sysdirect

import (
	client "github.com/i-Things/things/src/syssvr/client/user"
	server "github.com/i-Things/things/src/syssvr/internal/server/user"
)

func NewUser(config *Config) client.User {
	userSvc := getCtxSvc(config)
	return client.NewDirectUser(userSvc, server.NewUserServer(userSvc))
}
