package deviceGroup

import (
	"github.com/i-Things/things/src/dmsvr/internal/middleware"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
)

func InitMiddleware(svcCtx *svc.ServiceContext) {
	middleware.WithMiddlewares(middleware.DeviceDelete, svcCtx, DeviceDeleteHandle)
}
