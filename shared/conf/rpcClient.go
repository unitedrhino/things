package conf

import "github.com/zeromicro/go-zero/zrpc"

const (
	ClientModeGrpc   = "grpc"
	ClientModeDirect = "direct"
)

type RpcClientConf struct {
	Conf   zrpc.RpcClientConf `json:",optional"`
	Mode   string             `json:",default=grpc,options=grpc|direct"`
	Enable bool
}
