package otaFirmware

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.FirmwareUpdateReq) (resp *types.FirmwareResp, err error) {
	var firmwareUpdateReq dm.OtaFirmwareUpdateReq
	_ = copier.Copy(&firmwareUpdateReq, &req)
	firmwareUpdateReq.FirmwareUdi = &wrappers.StringValue{
		Value: req.FirmwareUdi,
	}
	logx.Infof("firmwareUpdateReq:%+v", &firmwareUpdateReq)
	update, err := l.svcCtx.OtaFirmwareM.OtaFirmwareUpdate(l.ctx, &firmwareUpdateReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.OtaFirmwareUpdate req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.FirmwareResp{FirmwareID: update.FirmwareID}, nil
}
