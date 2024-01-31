package middleware

import (
	operLog "gitee.com/i-Things/core/service/syssvr/client/log"
	"gitee.com/i-Things/core/shared/ctxs"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type SetupWareMiddleware struct {
	cfg    config.Config
	LogRpc operLog.Log
}

func NewSetupWareMiddleware(cfg config.Config, LogRpc operLog.Log) *SetupWareMiddleware {
	return &SetupWareMiddleware{cfg: cfg, LogRpc: LogRpc}
}

func (m *SetupWareMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.WithContext(r.Context()).Infof("%s.Lifecycle.Before", utils.FuncName())

		ctx2 := ctxs.SetMetaCtx(r.Context(), r.Header)
		r = r.WithContext(ctx2)

		next(w, r)

		logx.WithContext(r.Context()).Infof("%s.Lifecycle.After", utils.FuncName())
	}
}
