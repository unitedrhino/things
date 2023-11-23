package otafirmwaremanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryOTAFirmwareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB  *relationDB.ProductInfoRepo
	OfDB  *relationDB.OtaFirmwareRepo
	OffDB *relationDB.OtaFirmwareFileRepo
}

func NewQueryOTAFirmwareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryOTAFirmwareLogic {
	return &QueryOTAFirmwareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		OfDB:   relationDB.NewOtaFirmwareRepo(ctx),
		OffDB:  relationDB.NewOtaFirmwareFileRepo(ctx),
	}
}

// 查询升级包
func (l *QueryOTAFirmwareLogic) QueryOTAFirmware(in *dm.QueryOtaFirmwareReq) (*dm.QueryOtaFirmwareResp, error) {
	otaFirmware, err := l.OfDB.FindOneByFilter(l.ctx, relationDB.OtaFirmwareFilter{FirmwareID: in.FirmwareId})
	if err != nil {
		l.Errorf("%s.Query OTAFirmware err=%v", utils.FuncName(), err)
		return nil, err
	}
	otaFirmwareList, err := l.OffDB.FindByFilter(l.ctx, relationDB.OtaFirmwareFileFilter{FirmwareID: in.FirmwareId}, nil)
	if err != nil {
		l.Errorf("%s.Query OTAFirmwareFile err=%v", utils.FuncName(), err)
		return nil, err
	}
	var result *dm.QueryOtaFirmwareResp
	err = copier.Copy(&result, &otaFirmware)
	if err != nil {
		l.Errorf("%s.Copy OTAFirmware err=%v", utils.FuncName(), err)
		return nil, err
	}
	err = copier.Copy(&result.FirmwareFileList, &otaFirmwareList)
	if err != nil {
		l.Errorf("%s.Copy OTAFirmwareFile err=%v", utils.FuncName(), err)
		return nil, err
	}
	return result, nil
}
