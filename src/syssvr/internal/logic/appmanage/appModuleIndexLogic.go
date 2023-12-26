package appmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppModuleIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppModuleIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppModuleIndexLogic {
	return &AppModuleIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AppModuleIndexLogic) AppModuleIndex(in *sys.AppModuleIndexReq) (*sys.AppModuleIndexResp, error) {
	list, err := relationDB.NewAppModuleRepo(l.ctx).FindByFilter(l.ctx, relationDB.AppModuleFilter{AppCode: []string{in.Code}}, nil)
	if err != nil {
		return nil, err
	}
	var moduleCodes []string
	for _, v := range list {
		moduleCodes = append(moduleCodes, v.ModuleCode)
	}
	return &sys.AppModuleIndexResp{ModuleCodes: moduleCodes}, nil
}
