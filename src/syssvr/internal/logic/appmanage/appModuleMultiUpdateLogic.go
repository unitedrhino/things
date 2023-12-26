package appmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppModuleMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppModuleMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppModuleMultiUpdateLogic {
	return &AppModuleMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AppModuleMultiUpdateLogic) AppModuleMultiUpdate(in *sys.AppModuleMultiUpdateReq) (*sys.Response, error) {
	err := relationDB.NewAppModuleRepo(l.ctx).MultiUpdate(l.ctx, in.Code, in.ModuleCodes)
	return &sys.Response{}, err
}
