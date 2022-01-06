package dm

import (
	"context"
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/shared/utils"
	"github.com/go-things/things/src/dmsvr/dm"

	"github.com/go-things/things/src/webapi/internal/svc"
	"github.com/go-things/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetDeviceLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeviceLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetDeviceLogLogic {
	return GetDeviceLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDeviceLogLogic) GetDeviceLog(req types.GetDeviceLogReq) (*types.GetDeviceLogResp, error) {
	l.Infof("GetDeviceLog|req=%+v", req)
	resp,err := l.svcCtx.DmRpc.GetDeviceLog(l.ctx, &dm.GetDeviceLogReq{
		Method: req.Method,
		DeviceName: req.DeviceName,
		ProductID: req.ProductID,
		DataID: req.DataID,
		TimeStart: req.TimeStart,
		TimeEnd: req.TimeEnd,
		Limit: req.Limit,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceLog|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	info := make([]*types.DeviceData,0,len(resp.Data))
	for _,v := range resp.Data{
		info = append(info, &types.DeviceData{
			Timestamp:v.Timestamp,
			Method:v.Method,
			DataID:v.DataID,
			Payload:v.Payload,
		})
	}
	return &types.GetDeviceLogResp{
		Total: resp.Total,
		Data:info,
	}, nil
}
