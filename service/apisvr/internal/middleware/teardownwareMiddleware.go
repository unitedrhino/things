package middleware

import (
	operLog "gitee.com/unitedrhino/core/service/syssvr/client/log"
	"gitee.com/unitedrhino/things/service/apisvr/internal/config"
	"net/http"
	"sync"
)

type TeardownWareMiddleware struct {
	cfg    config.Config
	LogRpc operLog.Log
}

var respPool sync.Pool
var bufferSize = 1024

func init() {
	respPool.New = func() interface{} {
		return make([]byte, bufferSize)
	}
}

func NewTeardownWareMiddleware(cfg config.Config, LogRpc operLog.Log) *TeardownWareMiddleware {
	return &TeardownWareMiddleware{cfg: cfg, LogRpc: LogRpc}
}

func (m *TeardownWareMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//logx.WithContext(r.Context()).Infof("%s.Lifecycle.Before", utils.FuncName())

		next(w, r)

		//logx.WithContext(r.Context()).Infof("%s.Lifecycle.After", utils.FuncName())
	}
}
