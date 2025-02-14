package direct

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var (
	registerServers []func(server *rest.Server) error
)

func RegisterServer(run func(server *rest.Server) error) {
	registerServers = append(registerServers, run)
}

func InitServers(svr *rest.Server) *rest.Server {
	for _, r := range registerServers {
		err := r(svr)
		logx.Must(err)
	}
	return svr
}
