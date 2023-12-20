package appmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppInfoDeleteLogic {
	return &AppInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AppInfoDeleteLogic) AppInfoDelete(in *sys.WithIDCode) (*sys.Response, error) {
	err := relationDB.NewAppInfoRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.AppInfoFilter{Codes: []string{in.Code}, ID: in.Id})
	return &sys.Response{}, err
}
