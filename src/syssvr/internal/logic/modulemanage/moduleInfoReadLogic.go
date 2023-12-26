package modulemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleInfoReadLogic {
	return &ModuleInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleInfoReadLogic) ModuleInfoRead(in *sys.WithIDCode) (*sys.ModuleInfo, error) {
	ret, err := relationDB.NewModuleInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.ModuleInfoFilter{Codes: []string{in.Code}, ID: in.Id})
	return ToModuleInfoPb(ret), err
}
