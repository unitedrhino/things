package collect

import (
	"context"
	"github.com/i-Things/things/service/apisvr/internal/logic/things"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index() (resp *types.UserCollectDeviceInfo, err error) {
	ret, err := l.svcCtx.UserDevice.UserDeviceCollectIndex(l.ctx, &dm.Empty{})
	if err != nil {
		return nil, err
	}
	if len(ret.Devices) == 0 {
		return &types.UserCollectDeviceInfo{}, nil
	}
	var devs []*dm.DeviceCore
	for _, v := range ret.Devices {
		devs = append(devs, &dm.DeviceCore{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	ret2, err := l.svcCtx.DeviceM.DeviceInfoIndex(l.ctx, &dm.DeviceInfoIndexReq{Devices: devs})
	if err != nil {
		return nil, err
	}
	pis := make([]*types.DeviceInfo, 0, len(ret2.List))
	for _, v := range ret2.List {
		pi := things.InfoToApi(l.ctx, l.svcCtx, v, nil, nil, false)
		pis = append(pis, pi)
	}
	return &types.UserCollectDeviceInfo{Devices: pis}, nil
}
