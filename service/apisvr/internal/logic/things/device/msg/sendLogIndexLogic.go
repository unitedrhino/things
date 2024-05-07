package msg

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendLogIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendLogIndexLogic {
	return &SendLogIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *SendLogIndexLogic) SendLogIndex(req *types.DeviceMsgSendLogIndexReq) (resp *types.DeviceMsgSendLogIndexResp, err error) {
	dmResp, err := l.svcCtx.DeviceMsg.SendLogIndex(l.ctx, &dm.SendLogIndexReq{
		DeviceName: req.DeviceName,
		ProductID:  req.ProductID,
		TimeStart:  req.TimeStart,
		TimeEnd:    req.TimeEnd,
		Page:       logic.ToDiPageRpc(req.Page),
		UserID:     req.UserID,
		Actions:    req.Actions,
		ResultCode: req.ResultCode,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.HubLogIndex req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	info := make([]*types.DeviceMsgSendLogInfo, 0, len(dmResp.List))
	for _, v := range dmResp.List {
		info = append(info, &types.DeviceMsgSendLogInfo{
			Timestamp:  v.Timestamp,
			Account:    v.Account,
			UserID:     v.UserID,
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
			Action:     v.Action,
			DataID:     v.DataID,
			TraceID:    v.TraceID,
			Content:    v.Content,
			ResultCode: v.ResultCode,
		})
	}
	return &types.DeviceMsgSendLogIndexResp{List: info, Total: dmResp.Total}, err
}
