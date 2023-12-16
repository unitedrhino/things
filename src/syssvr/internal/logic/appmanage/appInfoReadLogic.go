package appmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppInfoReadLogic {
	return &AppInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AppInfoReadLogic) AppInfoRead(in *sys.ReqWithIDCode) (*sys.AppInfo, error) {
	ret, err := relationDB.NewAppInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.AppInfoFilter{Codes: []string{in.Code}, ID: in.Id})
	return ToAppInfoPb(ret), err
}
