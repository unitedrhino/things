package otamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

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
	//todo debug
	//if err := ctxs.IsRoot(l.ctx); err != nil {
	//	return nil, err
	//}
	l.ctx = ctxs.WithRoot(l.ctx)
	err := relationDB.NewOtaModuleInfoRepo(l.ctx).Delete(l.ctx, in.Id)
	return &dm.Empty{}, err
}
