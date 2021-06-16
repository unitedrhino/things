package msvc

import (
	"gitee.com/godLei6/things/src/dmsvr/device/msgquque/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
