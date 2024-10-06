package share

import (
	"context"
	"gitee.com/i-Things/things/service/apisvr/internal/logic/things"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"gitee.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.UserDeviceShareReadReq) (resp *types.UserDeviceShareInfo, err error) {
	ret, err := l.svcCtx.UserDevice.UserDeviceShareRead(l.ctx, &dm.UserDeviceShareReadReq{
		Id:     req.ID,
		Device: things.ToDmDeviceCorePb(req.Device),
	})
	if err != nil {
		return nil, err
	}
	return ToShareTypes(ret), nil
}
