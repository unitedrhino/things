package userdirect

import (
	"github.com/i-Things/things/src/usersvr/internal/config"
	"github.com/i-Things/things/src/usersvr/internal/server"
	"github.com/i-Things/things/src/usersvr/internal/svc"
	"github.com/i-Things/things/src/usersvr/userclient"
)

type Config = config.Config

var ctxSvc *svc.ServiceContext

func NewUser(config Config) userclient.User {
	userSvc := svc.NewServiceContext(config)
	return userclient.NewDirectUser(userSvc, server.NewUserServer(userSvc))
}
