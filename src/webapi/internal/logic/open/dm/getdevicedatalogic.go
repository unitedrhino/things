package dm

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/dm"

	"github.com/i-Things/things/src/webapi/internal/svc"
	"github.com/i-Things/things/src/webapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeviceDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeviceDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetDeviceDataLogic {
	return GetDeviceDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDeviceDataLogic) GetDeviceData(req types.GetDeviceDataReq) (*types.GetDeviceDataResp, error) {
	l.Infof("GetDeviceData|req=%+v", req)
	resp, err := l.svcCtx.DmRpc.GetDeviceData(l.ctx, &dm.GetDeviceDataReq{
		Method:     req.Method,
		DeviceName: req.DeviceName,
		ProductID:  req.ProductID,
		DataID:     req.DataID,
		TimeStart:  req.TimeStart,
		TimeEnd:    req.TimeEnd,
		Limit:      req.Limit,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceData|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	info := make([]*types.DeviceData, 0, len(resp.List))
	for _, v := range resp.List {
		info = append(info, &types.DeviceData{
			Timestamp: v.Timestamp,
			Type:      v.Type,
			DataID:    v.DataID,
			GetValue:  v.GetValue,
			SendValue: v.SendValue,
		})
	}
	return &types.GetDeviceDataResp{
		Total: resp.Total,
		List:  info,
	}, nil
}
