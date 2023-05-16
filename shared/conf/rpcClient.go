package conf

import "github.com/zeromicro/go-zero/zrpc"

const (
	ClientModeGrpc   = "grpc"
	ClientModeDirect = "direct"
)

type RpcClientConf struct {
	Conf zrpc.RpcClientConf `json:",optional"`
	ModeConf
}

type ModeConf struct {
	Mode     string `json:",default=direct,options=grpc|direct"`
	RunProxy bool   `json:",default=false"` //是否开启grpc服务
	Enable   bool   `json:",default=true"`
}
