package dm

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeviceDescribeLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeviceDescribeLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetDeviceDescribeLogLogic {
	return GetDeviceDescribeLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDeviceDescribeLogLogic) GetDeviceDescribeLog(req types.GetDeviceDescribeLogReq) (*types.GetDeviceDescribeLogResp, error) {
	l.Infof("GetDeviceDescribeLog|req=%+v", req)
	resp, err := l.svcCtx.DmRpc.GetDeviceDescribeLog(l.ctx, &dm.GetDeviceDescribeLogReq{
		DeviceName: req.DeviceName,
		ProductID:  req.ProductID,
		TimeStart:  req.TimeStart,
		TimeEnd:    req.TimeEnd,
		Limit:      req.Limit,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceDescribeLog|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	info := make([]*types.DeviceDescribeLog, 0, len(resp.List))
	for _, v := range resp.List {
		info = append(info, &types.DeviceDescribeLog{
			Timestamp:  v.Timestamp,
			Action:     v.Action,
			RequestID:  v.RequestID,
			TranceID:   v.TranceID,
			Topic:      v.Topic,
			Content:    v.Content,
			ResultType: v.ResultType,
		})
	}
	return &types.GetDeviceDescribeLogResp{List: info}, err
}
