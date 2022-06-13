package logic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/spf13/cast"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ManageFirmwareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewManageFirmwareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManageFirmwareLogic {
	return &ManageFirmwareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ManageFirmwareLogic) AddFirmware(in *dm.ManageFirmwareReq) (*dm.FirmwareInfo, error) {
	return nil, nil
}

func (l *ManageFirmwareLogic) ModifyFirmware(in *dm.ManageFirmwareReq) (*dm.FirmwareInfo, error) {
	return nil, nil
}

func (l *ManageFirmwareLogic) DelFirmware(in *dm.ManageFirmwareReq) (*dm.FirmwareInfo, error) {
	return nil, nil
}

// 管理产品的固件
func (l *ManageFirmwareLogic) ManageFirmware(in *dm.ManageFirmwareReq) (*dm.FirmwareInfo, error) {
	l.Infof("[%s]opt=%d|info=%+v", utils.FuncName(), in.Opt, in.Info)
	switch in.Opt {
	case def.OPT_ADD:
		if in.Info == nil {
			return nil, errors.Parameter.WithMsg("add opt need info")
		}
		return l.AddFirmware(in)
	case def.OPT_MODIFY:
		return l.ModifyFirmware(in)
	case def.OPT_DEL:
		return l.DelFirmware(in)
	default:
		return nil, errors.Parameter.AddDetail("not support opt:" + cast.ToString(in.Opt))
	}
}
