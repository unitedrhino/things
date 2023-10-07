package config

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	Database conf.Database
	Event    conf.EventConf
	zrpc.RpcServerConf
}
