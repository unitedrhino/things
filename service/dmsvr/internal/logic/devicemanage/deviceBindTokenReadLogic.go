package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceBindTokenReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceBindTokenReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceBindTokenReadLogic {
	return &DeviceBindTokenReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceBindTokenReadLogic) DeviceBindTokenRead(in *dm.DeviceBindTokenReadReq) (*dm.DeviceBindTokenInfo, error) {
	tk, err := l.svcCtx.DeviceBindToken.GetData(l.ctx, in.Token)
	if err != nil {
		return nil, err
	}
	return &dm.DeviceBindTokenInfo{Token: tk.Token, Status: tk.Status}, nil
}
