package firmwaremanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type FirmwareInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFirmwareInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FirmwareInfoUpdateLogic {
	return &FirmwareInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FirmwareInfoUpdateLogic) ChangeDevice(old *relationDB.DmOtaFirmware, data *dm.FirmwareInfo) {
	if data.Name != "" {
		old.Name = data.Name
	}
	if data.Desc != nil {
		old.Desc = data.Desc.Value
	}
	//if data.ExtData != nil {
	//	old.Extra = sql.NullString{
	//		String: data.ExtData.Value,
	//		Valid:  true,
	//	}
	//}
}

func (l *FirmwareInfoUpdateLogic) FirmwareInfoUpdate(in *dm.FirmwareInfo) (*dm.OtaCommonResp, error) {
	var fDB = relationDB.NewOtaFirmwareRepo(l.ctx)
	di, err := fDB.FindOne(l.ctx, in.FirmwareID)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.NotFind.AddDetailf("not find firmware|firmwareID=%s|name=%s",
				in.FirmwareID, in.Name)
		}
		return nil, errors.Database.AddDetail(err)
	}
	l.ChangeDevice(di, in)
	err = fDB.Update(l.ctx, di)
	if err != nil {
		l.Errorf("DeviceInfoUpdate.DeviceInfo.Update err=%+v", err)
		return nil, errors.System.AddDetail(err)
	}
	return &dm.OtaCommonResp{}, nil
}
