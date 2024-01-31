package accessmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiInfoIndexLogic {
	return &ApiInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiInfoIndexLogic) ApiInfoIndex(in *sys.ApiInfoIndexReq) (*sys.ApiInfoIndexResp, error) {
	f := relationDB.ApiInfoFilter{
		ApiIDs:       in.ApiIDs,
		Route:        in.Route,
		Method:       in.Method,
		Name:         in.Name,
		AccessCode:   in.AccessCode,
		IsAuthTenant: in.IsAuthTenant,
	}
	pos, err := relationDB.NewApiInfoRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	total, err := relationDB.NewApiInfoRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	var ais []*sys.ApiInfo
	for _, v := range pos {
		ais = append(ais, ToApiInfoPb(v))
	}
	return &sys.ApiInfoIndexResp{
		List:  ais,
		Total: total,
	}, nil
}
