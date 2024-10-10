package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/errors"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoMultiBindLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceInfoMultiBindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoMultiBindLogic {
	return &DeviceInfoMultiBindLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceInfoMultiBindLogic) DeviceInfoMultiBind(in *dm.DeviceInfoMultiBindReq) (*dm.DeviceInfoMultiBindResp, error) {
	bind := NewDeviceInfoBindLogic(l.ctx, l.svcCtx)
	var errs []*dm.DeviceError
	for _, dev := range in.Devices {
		_, err := bind.DeviceInfoBind(&dm.DeviceInfoBindReq{Device: dev, AreaID: in.AreaID})
		if err == nil {
			continue
		}
		er := errors.Fmt(err)
		errs = append(errs, &dm.DeviceError{
			ProductID:  dev.ProductID,
			DeviceName: dev.DeviceName,
			Msg:        er.GetMsg(),
			Code:       er.GetCode(),
		})
	}

	return &dm.DeviceInfoMultiBindResp{Errs: errs}, nil
}
