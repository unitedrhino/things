package apimanagelogic

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
	AiDB *relationDB.ApiInfoRepo
}

func NewApiInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiInfoIndexLogic {
	return &ApiInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewApiInfoRepo(ctx),
	}
}

func (l *ApiInfoIndexLogic) ApiInfoIndex(in *sys.ApiInfoIndexReq) (*sys.ApiInfoIndexResp, error) {
	f := relationDB.ApiInfoFilter{
		Route:   in.Route,
		Method:  in.Method,
		Group:   in.Group,
		Name:    in.Name,
		AppCode: in.AppCode,
	}
	resp, err := l.AiDB.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	total, err := l.AiDB.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	info := make([]*sys.ApiInfo, 0, len(resp))
	for _, v := range resp {
		info = append(info, &sys.ApiInfo{
			Id:           v.ID,
			Route:        v.Route,
			Method:       v.Method,
			Group:        v.Group,
			Name:         v.Name,
			BusinessType: v.BusinessType,
			AppCode:      v.AppCode,
		})
	}

	return &sys.ApiInfoIndexResp{
		Total: total,
		List:  info,
	}, nil
}
