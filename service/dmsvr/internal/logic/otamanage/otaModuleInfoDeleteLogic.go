package otamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaModuleInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaModuleInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaModuleInfoDeleteLogic {
	return &OtaModuleInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OtaModuleInfoDeleteLogic) OtaModuleInfoDelete(in *dm.WithID) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	l.ctx = ctxs.WithRoot(l.ctx)
	err := relationDB.NewOtaModuleInfoRepo(l.ctx).Delete(l.ctx, in.Id)
	return &dm.Empty{}, err
}
