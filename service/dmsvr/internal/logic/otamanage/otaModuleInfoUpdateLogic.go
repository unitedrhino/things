package otamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaModuleInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaModuleInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaModuleInfoUpdateLogic {
	return &OtaModuleInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OtaModuleInfoUpdateLogic) OtaModuleInfoUpdate(in *dm.OtaModuleInfo) (*dm.Empty, error) {
	//todo debug
	//if err := ctxs.IsRoot(l.ctx); err != nil {
	//	return nil, err
	//}
	l.ctx = ctxs.WithRoot(l.ctx)
	old, err := relationDB.NewOtaModuleInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	old.Name = in.Name
	old.Desc = in.Desc
	err = relationDB.NewOtaModuleInfoRepo(l.ctx).Update(l.ctx, old)
	return &dm.Empty{}, err
}
