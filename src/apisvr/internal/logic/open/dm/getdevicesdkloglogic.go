package dm

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeviceSDKLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeviceSDKLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceSDKLogLogic {
	return &GetDeviceSDKLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDeviceSDKLogLogic) GetDeviceSDKLog(req *types.GetDeviceSDKLogReq) (*types.GetDeviceSDKLogResp, error) {
	dmReq := &dm.GetDeviceSDKLogReq{
		DeviceName: req.DeviceName,
		ProductID:  req.ProductID, //产品id
		TimeStart:  req.TimeStart,
		TimeEnd:    req.TimeEnd,
		Limit:      req.Limit,
	}
	if req.Page == 0 {
		req.Page = 1
	}
	dmReq.Page = &dm.PageInfo{
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	resp, err := l.svcCtx.DmRpc.GetDeviceSDKLog(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceSDKLog|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	info := make([]*types.DeviceSDKLog, 0, len(resp.List))
	for _, v := range resp.List {
		info = append(info, &types.DeviceSDKLog{
			Timestamp: v.Timestamp,
			Loglevel:  v.Loglevel,
			Content:   v.Content,
		})
	}
	return &types.GetDeviceSDKLogResp{List: info}, err
}
