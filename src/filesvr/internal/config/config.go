package config

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	Rpc  zrpc.RpcServerConf
	Rest rest.RestConf
	OSS  conf.OSSConf
}
