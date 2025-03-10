package msg

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

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
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *EventLogIndexLogic) EventLogIndex(req *types.DeviceMsgEventLogIndexReq) (resp *types.DeviceMsgEventLogIndexResp, err error) {
	dmResp, err := l.svcCtx.DeviceMsg.EventLogIndex(l.ctx, &dm.EventLogIndexReq{
		DeviceNames: req.DeviceNames,
		ProductID:   req.ProductID,
		DataID:      req.DataID,
		TimeStart:   req.TimeStart,
		TimeEnd:     req.TimeEnd,
		Types:       req.Types,
		Page:        logic.ToDmPageRpc(req.Page),
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.GetDeviceData req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	info := make([]*types.DeviceMsgEventLogInfo, 0, len(dmResp.List))
	for _, v := range dmResp.List {
		info = append(info, &types.DeviceMsgEventLogInfo{
			Timestamp: v.Timestamp,
			DataID:    v.DataID,
			Params:    v.Params,
			Type:      v.Type,
		})
	}
	return &types.DeviceMsgEventLogIndexResp{
		PageResp: logic.ToPageResp(req.Page, dmResp.Total),
		List:     info,
	}, nil
}
