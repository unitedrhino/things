package msg

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/disvr/pb/di"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyLogIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPropertyLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyLogIndexLogic {
	return &PropertyLogIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyLogIndexLogic) PropertyLogIndex(req *types.DeviceMsgPropertyLogIndexReq) (resp *types.DeviceMsgPropertyIndexResp, err error) {
	dmResp, err := l.svcCtx.DeviceMsg.PropertyLogIndex(l.ctx, &di.PropertyLogIndexReq{
		DeviceNames: req.DeviceNames,
		ProductID:   req.ProductID,
		DataID:      req.DataID,
		TimeStart:   req.TimeStart,
		TimeEnd:     req.TimeEnd,
		Interval:    req.Interval,
		ArgFunc:     req.ArgFunc,
		Fill:        req.Fill,
		Order:       req.Order,
		Page:        logic.ToDiPageRpc(req.Page),
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.GetDeviceData req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	info := make([]*types.DeviceMsgPropertyIndex, 0, len(dmResp.List))
	for _, v := range dmResp.List {
		info = append(info, &types.DeviceMsgPropertyIndex{
			Timestamp: v.Timestamp,
			DataID:    v.DataID,
			Value:     v.Value,
		})
	}
	return &types.DeviceMsgPropertyIndexResp{
		Total: dmResp.Total,
		List:  info,
	}, nil
}
