package userdirect

import (
	"github.com/i-Things/things/src/usersvr/internal/config"
	"github.com/i-Things/things/src/usersvr/internal/server"
	"github.com/i-Things/things/src/usersvr/internal/svc"
	"github.com/i-Things/things/src/usersvr/user"
)

type Config = config.Config

var ctxSvc *svc.ServiceContext

func NewUser(config Config) user.User {
	userSvc := svc.NewServiceContext(config)
	return user.NewDirectUser(userSvc, server.NewUserServer(userSvc))
}
