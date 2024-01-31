package otafirmwaremanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
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
	otaFirmware, err := l.OfDB.FindOneByFilter(l.ctx, relationDB.OtaFirmwareFilter{FirmwareID: in.FirmwareId})
	if err != nil {
		return nil, err
	}
	//更新相关字段
	otaFirmware.Desc = in.FirmwareDesc
	otaFirmware.Name = in.FirmwareName
	otaFirmware.Extra = in.FirmwareUdi.Value
	err = l.OfDB.Update(l.ctx, otaFirmware)
	if err != nil {
		l.Errorf("%s.Update err=%v", utils.FuncName(), err)
		return nil, errors.System.AddDetail(err)
	}
	return &dm.OtaFirmwareResp{FirmwareID: otaFirmware.ID}, nil
}
