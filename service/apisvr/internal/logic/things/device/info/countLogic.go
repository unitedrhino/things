package info

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

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
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CountLogic) Count(req *types.DeviceCountReq) (resp *types.DeviceCountResp, err error) {
	diReq := &dm.DeviceCountReq{
		CountTypes: req.CountTypes,
		RangeType:  req.RangeType,
		RangeIDs:   req.RangeIDs,
	}
	diResp, err := l.svcCtx.DeviceM.DeviceCount(l.ctx, diReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.DeviceInfoCount req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	var deviceCountInfos []*types.DeviceCountInfo
	for _, v := range diResp.List {
		deviceCountInfos = append(deviceCountInfos, &types.DeviceCountInfo{RangeID: v.RangeID, Count: v.Count})
	}
	return &types.DeviceCountResp{
		List: deviceCountInfos,
	}, nil
}
