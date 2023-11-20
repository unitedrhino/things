package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CountLogic {
	return &CountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CountLogic) Count(req *types.DeviceCountReq) (resp *types.DeviceCountResp, err error) {
	diReq := &dm.DeviceInfoCountReq{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	diResp, err := l.svcCtx.DeviceM.DeviceInfoCount(l.ctx, diReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.DeviceInfoCount req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	dtReq := &dm.DeviceTypeCountReq{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	dtResp, err := l.svcCtx.DeviceM.DeviceTypeCount(l.ctx, dtReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.DeviceTypeCount req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}

	return &types.DeviceCountResp{
		DeviceInfoCount: types.DeviceInfoCount{
			Online:   diResp.Online,
			Offline:  diResp.Offline,
			Inactive: diResp.Inactive,
			Unknown:  diResp.Unknown,
		},
		DeviceTypeCount: types.DeviceTypeCount{
			Device:  dtResp.Device,
			Gateway: dtResp.Gateway,
			Subset:  dtResp.Subset,
			Unknown: dtResp.Unknown,
		},
	}, nil
}
