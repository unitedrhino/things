package otafirmwaremanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaFirmwareUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
	OfDB *relationDB.OtaFirmwareRepo
}

func NewOtaFirmwareUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaFirmwareUpdateLogic {
	return &OtaFirmwareUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		OfDB:   relationDB.NewOtaFirmwareRepo(ctx),
	}
}

// 修改升级包
func (l *OtaFirmwareUpdateLogic) OtaFirmwareUpdate(in *dm.OtaFirmwareUpdateReq) (*dm.OtaFirmwareResp, error) {
	var otaFirmware relationDB.DmOtaFirmware
	copier.Copy(&otaFirmware, in)
	logx.Infof("otaFirmware:%+v", otaFirmware)
	err := l.OfDB.Update(l.ctx, &otaFirmware)
	if err != nil {
		l.Errorf("%s.Update err=%v", utils.FuncName(), err)
		return nil, errors.System.AddDetail(err)
	}
	return &dm.OtaFirmwareResp{FirmwareID: otaFirmware.ID}, nil
}
