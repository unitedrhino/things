package conf

import "github.com/tal-tech/go-zero/zrpc"

type RpcClientConf struct {
	Conf   zrpc.RpcClientConf `json:",optional"`
	Enable bool
}
