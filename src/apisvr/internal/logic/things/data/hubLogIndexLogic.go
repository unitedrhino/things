package data

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

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

func (l *HubLogIndexLogic) HubLogIndex(req *types.DataHubLogIndexReq) (resp *types.DataHubLogIndexResp, err error) {
	dmResp, err := l.svcCtx.DmRpc.DataHubLogIndex(l.ctx, &dm.DataHubLogIndexReq{
		DeviceName: req.DeviceName,
		ProductID:  req.ProductID,
		TimeStart:  req.TimeStart,
		TimeEnd:    req.TimeEnd,
		Page: &dm.PageInfo{
			Page: req.Page.Page,
			Size: req.Page.Size,
		},
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceDescribeLog|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	info := make([]*types.DataHubLogIndex, 0, len(dmResp.List))
	for _, v := range dmResp.List {
		info = append(info, &types.DataHubLogIndex{
			Timestamp:  v.Timestamp,
			Action:     v.Action,
			RequestID:  v.RequestID,
			TranceID:   v.TranceID,
			Topic:      v.Topic,
			Content:    v.Content,
			ResultType: v.ResultType,
		})
	}
	return &types.DataHubLogIndexResp{List: info}, err
}
