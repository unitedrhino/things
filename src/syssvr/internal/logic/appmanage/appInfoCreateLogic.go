package appmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppInfoCreateLogic {
	return &AppInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AppInfoCreateLogic) AppInfoCreate(in *sys.AppInfo) (*sys.Response, error) {
	in.Id = 0
	err := relationDB.NewAppInfoRepo(l.ctx).Insert(l.ctx, ToAppInfoPo(in))
	return &sys.Response{}, err
}
