package otamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaFirmwareInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
	OfDB *relationDB.OtaFirmwareInfoRepo
}

func NewOtaFirmwareInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaFirmwareInfoUpdateLogic {
	return &OtaFirmwareInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		OfDB:   relationDB.NewOtaFirmwareInfoRepo(ctx),
	}
}

// 修改升级包
func (l *OtaFirmwareInfoUpdateLogic) OtaFirmwareInfoUpdate(in *dm.OtaFirmwareInfoUpdateReq) (*dm.WithID, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	l.ctx = ctxs.WithRoot(l.ctx)
	otaFirmware, err := l.OfDB.FindOneByFilter(l.ctx, relationDB.OtaFirmwareInfoFilter{ID: in.Id})
	if err != nil {
		return nil, err
	}
	//更新相关字段
	otaFirmware.Desc = in.Desc
	otaFirmware.Name = in.Name
	otaFirmware.Extra = in.Extra.GetValue()
	err = l.OfDB.Update(l.ctx, otaFirmware)
	if err != nil {
		l.Errorf("%s.Update err=%v", utils.FuncName(), err)
		return nil, err
	}
	return &dm.WithID{Id: otaFirmware.ID}, nil
}
