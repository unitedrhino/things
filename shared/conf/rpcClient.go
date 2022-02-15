package conf

import "github.com/zeromicro/go-zero/zrpc"

type RpcClientConf struct {
	Conf   zrpc.RpcClientConf `json:",optional"`
	Enable bool
}
