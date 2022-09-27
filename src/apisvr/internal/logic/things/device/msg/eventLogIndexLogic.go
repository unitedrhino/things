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

type EventLogIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEventLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EventLogIndexLogic {
	return &EventLogIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EventLogIndexLogic) EventLogIndex(req *types.DeviceMsgEventLogIndexReq) (resp *types.DeviceMsgEventIndexResp, err error) {
	dmResp, err := l.svcCtx.DeviceMsg.EventLogIndex(l.ctx, &di.EventLogIndexReq{
		DeviceNames: req.DeviceNames,
		ProductID:   req.ProductID,
		DataID:      req.DataID,
		TimeStart:   req.TimeStart,
		TimeEnd:     req.TimeEnd,
		Types:       req.Types,
		Page:        logic.ToDiPageRpc(req.Page),
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.GetDeviceData req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	info := make([]*types.DeviceMsgEventIndex, 0, len(dmResp.List))
	for _, v := range dmResp.List {
		info = append(info, &types.DeviceMsgEventIndex{
			Timestamp: v.Timestamp,
			DataID:    v.DataID,
			Params:    v.Params,
			Type:      v.Type,
		})
	}
	return &types.DeviceMsgEventIndexResp{
		Total: dmResp.Total,
		List:  info,
	}, nil
}
