package msg

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/logic"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

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
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HubLogIndexLogic) HubLogIndex(req *types.DeviceMsgHubLogIndexReq) (resp *types.DeviceMsgHubLogIndexResp, err error) {
	dmResp, err := l.svcCtx.DeviceMsg.HubLogIndex(l.ctx, &dm.HubLogIndexReq{
		DeviceName: req.DeviceName,
		ProductID:  req.ProductID,
		TimeStart:  req.TimeStart,
		TimeEnd:    req.TimeEnd,
		Page:       logic.ToDiPageRpc(req.Page),
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
	info := make([]*types.DeviceMsgHubLogIndex, 0, len(dmResp.List))
	for _, v := range dmResp.List {
		info = append(info, &types.DeviceMsgHubLogIndex{
			Timestamp:  v.Timestamp,
			Action:     v.Action,
			RequestID:  v.RequestID,
			TranceID:   v.TranceID,
			Topic:      v.Topic,
			Content:    v.Content,
			ResultType: v.ResultType,
		})
	}
	return &types.DeviceMsgHubLogIndexResp{List: info, Total: dmResp.Total}, err
}
