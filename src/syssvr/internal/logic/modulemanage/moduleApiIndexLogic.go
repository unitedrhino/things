package modulemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleApiIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleApiIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleApiIndexLogic {
	return &ModuleApiIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleApiIndexLogic) ModuleApiIndex(in *sys.ApiInfoIndexReq) (*sys.ApiInfoIndexResp, error) {
	f := relationDB.ApiInfoFilter{
		ApiIDs:     in.ApiIDs,
		Route:      in.Route,
		Method:     in.Method,
		Group:      in.Group,
		Name:       in.Name,
		ModuleCode: in.ModuleCode,
		IsNeedAuth: in.IsNeedAuth,
	}
	resp, err := relationDB.NewApiInfoRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	total, err := relationDB.NewApiInfoRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	info := make([]*sys.ApiInfo, 0, len(resp))
	for _, v := range resp {
		info = append(info, ToApiInfoPb(v))
	}

	return &sys.ApiInfoIndexResp{
		Total: total,
		List:  info,
	}, nil
}
