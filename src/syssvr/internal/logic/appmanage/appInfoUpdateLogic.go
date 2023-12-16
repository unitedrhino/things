package appmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppInfoUpdateLogic {
	return &AppInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AppInfoUpdateLogic) AppInfoUpdate(in *sys.AppInfo) (*sys.Response, error) {
	err := relationDB.NewAppInfoRepo(l.ctx).Update(l.ctx, ToAppInfoPo(in))
	return &sys.Response{}, err
}
