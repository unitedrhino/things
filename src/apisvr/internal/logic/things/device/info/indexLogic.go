package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/logic/things/device"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.DeviceInfoIndexReq) (resp *types.DeviceInfoIndexResp, err error) {
	dmReq := &dm.DeviceInfoIndexReq{
		ProductID:  req.ProductID, //产品id
		DeviceName: req.DeviceName,
		Tags:       logic.ToTagsMap(req.Tags),
		Page:       logic.ToDmPageRpc(req.Page),
		Range:      req.Range,
		Position:   logic.ToDmPointRpc(req.Position),
	}
	dmResp, err := l.svcCtx.DeviceM.DeviceInfoIndex(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.GetDeviceInfo req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	pis := make([]*types.DeviceInfo, 0, len(dmResp.List))
	for _, v := range dmResp.List {
		pi := device.InfoToApi(l.ctx, l.svcCtx, v, req.WithProperty)
		pis = append(pis, pi)
	}
	return &types.DeviceInfoIndexResp{
		Total: dmResp.Total,
		List:  pis,
		Num:   int64(len(pis)),
	}, nil
}
