package otafirmwaremanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/shared/utils/cast"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteOTAFirmwareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
	OfDB *relationDB.OtaFirmwareRepo
}

func NewDeleteOTAFirmwareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOTAFirmwareLogic {
	return &DeleteOTAFirmwareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		OfDB:   relationDB.NewOtaFirmwareRepo(ctx),
	}
}

// 删除升级包
func (l *DeleteOTAFirmwareLogic) DeleteOTAFirmware(in *dm.DeleteOtaFirmwareReq) (*dm.Response, error) {
	_, err := l.OfDB.FindOneByFilter(l.ctx, relationDB.OtaFirmwareFilter{FirmwareID: in.FirmwareId})
	if errors.Cmp(err, errors.NotFind) {
		l.Errorf("not find firmware id:" + cast.ToString(in.FirmwareId))
		return nil, err
	} else if err != nil {
		return nil, err
	}
	err = l.OfDB.Delete(l.ctx, in.FirmwareId)
	if err != nil {
		l.Errorf("%s.DeleteOTAFirmware err=%v", utils.FuncName(), err)
		return nil, errors.System.AddDetail(err)
	}
	return &dm.Response{}, nil
}
