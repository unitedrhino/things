package middleware

import (
	"github.com/i-Things/things/shared/middlewares"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
)

const (
	DeviceDelete = "deviceDelete"
)

type (
	dmMiddleware func(svcCtx *svc.ServiceContext, next middlewares.HandleFunc) middlewares.HandleFunc
)

func WithMiddlewares(name string, svcCtx *svc.ServiceContext, mds ...dmMiddleware) {
	for _, md := range mds {
		middlewares.WithMiddlewares(name, func(next middlewares.HandleFunc) middlewares.HandleFunc {
			return md(svcCtx, next)
		})
	}
}
