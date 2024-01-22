package job

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeviceIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceIndexLogic {
	return &DeviceIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeviceIndexLogic) DeviceIndex(req *types.OtaJobByDeviceIndexReq) (resp *types.OtaJobInfoIndexResp, err error) {
	// todo: add your logic here and delete this line

	return
}
