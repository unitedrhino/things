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

type HubLogIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHubLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HubLogIndexLogic {
	return &HubLogIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *HubLogIndexLogic) HubLogIndex(req *types.DeviceMsgHubLogIndexReq) (resp *types.DeviceMsgHubLogIndexResp, err error) {
	dmResp, err := l.svcCtx.DeviceMsg.HubLogIndex(l.ctx, &dm.HubLogIndexReq{
		DeviceName: req.DeviceName,
		ProductID:  req.ProductID,
		TimeStart:  req.TimeStart,
		TimeEnd:    req.TimeEnd,
		Page:       logic.ToDmPageRpc(req.Page),
		Actions:    req.Actions,
		Topics:     req.Topics,
		Content:    req.Content,
		RequestID:  req.RequestID,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.HubLogIndex req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	info := make([]*types.DeviceMsgHubLogInfo, 0, len(dmResp.List))
	for _, v := range dmResp.List {
		info = append(info, &types.DeviceMsgHubLogInfo{
			Timestamp:   v.Timestamp,
			Action:      v.Action,
			RequestID:   v.RequestID,
			TraceID:     v.TraceID,
			Topic:       v.Topic,
			Content:     v.Content,
			ResultCode:  v.ResultCode,
			RespPayload: v.RespPayload,
		})
	}
	return &types.DeviceMsgHubLogIndexResp{List: info, PageResp: logic.ToPageResp(req.Page, dmResp.Total)}, err
}
