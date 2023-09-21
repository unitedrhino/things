package firmware

import (
	"context"

	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeviceInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoReadLogic {
	return &DeviceInfoReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeviceInfoReadLogic) DeviceInfoRead(req *types.OtaFirmwareDeviceInfoReq) (resp *types.OtaFirmwareDeviceInfoResp, err error) {
	dmResp, err := l.svcCtx.FirmwareM.OtaFirmwareDeviceInfo(l.ctx, &dm.OtaFirmwareDeviceInfoReq{FirmwareID: req.FirmwareID})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}

	return &types.OtaFirmwareDeviceInfoResp{
		Versions: dmResp.GetVersions(),
	}, nil
}
