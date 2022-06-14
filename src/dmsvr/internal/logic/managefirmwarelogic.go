package logic

import (
	"context"
	"database/sql"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/spf13/cast"
	"time"

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
	_, err := l.svcCtx.FirmwareInfo.FindOneByProductIDVersion(l.ctx, in.Info.ProductID, in.Info.Version)
	if err == nil {
		return nil, errors.Duplicate.WithMsg("产品固件版本已有")
	}
	if err != mysql.ErrNotFound {
		l.Errorf("AddFirmware|FindOneByProductIDVersion|err=%v", err)
		return nil, errors.Database.AddDetail(err.Error())
	}
	firmware := mysql.ProductFirmware{
		ProductID:   in.Info.ProductID,
		Version:     in.Info.Version,
		CreatedTime: time.Now(),
		Name:        in.Info.Name,
		Description: in.Info.Description,
		Size:        in.Info.Size,
		Dir:         in.Info.Dir,
	}
	_, err = l.svcCtx.FirmwareInfo.Insert(l.ctx, &firmware)
	if err != nil {
		l.Errorf("[%s]Insert|err=%+v", err)
		return nil, errors.System.AddDetail(err.Error())
	}
	return in.Info, nil
}

func (l *ManageFirmwareLogic) ModifyFirmware(in *dm.ManageFirmwareReq) (*dm.FirmwareInfo, error) {
	oldFirmWare, err := l.svcCtx.FirmwareInfo.FindOneByProductIDVersion(l.ctx, in.Info.ProductID, in.Info.Version)
	if err != nil {
		if err != mysql.ErrNotFound {
			l.Errorf("AddFirmware|FindOneByProductIDVersion|err=%v", err)
			return nil, errors.Database.AddDetail(err.Error())
		}
		return nil, errors.NotFind
	}
	oldFirmWare.Name = in.Info.Name
	oldFirmWare.Description = in.Info.Description
	oldFirmWare.UpdatedTime = sql.NullTime{Valid: true, Time: time.Now()}
	err = l.svcCtx.FirmwareInfo.Update(l.ctx, oldFirmWare)
	if err != nil {
		return nil, errors.Database.AddDetail(err.Error())
	}
	return ToFirmwareInfo(oldFirmWare), nil
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
