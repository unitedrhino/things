package task

import (
	"context"

	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/apisvr/internal/logic"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeviceIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceIndexLogic {
	return &DeviceIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeviceIndexLogic) DeviceIndex(req *types.OtaTaskDeviceIndexReq) (resp *types.OtaTaskDeviceIndexResp, err error) {
	otaResp, err := l.svcCtx.OtaTaskM.OtaTaskDeviceIndex(l.ctx, &dm.OtaTaskDeviceIndexReq{
		TaskUid:    req.TaskUid,
		DeviceName: req.DeviceName,
		Status:     req.Status,
		Page:       logic.ToOtaPageRpc(req.Page),
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	pis := make([]*types.OtaTaskDeviceInfo, 0, len(otaResp.List))
	for _, v := range otaResp.List {
		pi := otaTaskDeviceInfoToApi(v)
		pis = append(pis, pi)
	}
	return &types.OtaTaskDeviceIndexResp{
		List:  pis,
		Total: otaResp.Total,
	}, nil
}
