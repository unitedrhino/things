package apilogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.ApiInfoRepo
}

func NewApiIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiIndexLogic {
	return &ApiIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewApiInfoRepo(ctx),
	}
}

func (l *ApiIndexLogic) ApiIndex(in *sys.ApiIndexReq) (*sys.ApiIndexResp, error) {
	f := relationDB.ApiInfoFilter{
		Route:  in.Route,
		Method: in.Method,
		Group:  in.Group,
		Name:   in.Name,
	}
	resp, err := l.AiDB.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	total, err := l.AiDB.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	info := make([]*sys.ApiData, 0, len(resp))
	for _, v := range resp {
		info = append(info, &sys.ApiData{
			Id:           v.ID,
			Route:        v.Route,
			Method:       v.Method,
			Group:        v.Group,
			Name:         v.Name,
			BusinessType: v.BusinessType,
		})
	}

	return &sys.ApiIndexResp{
		Total: total,
		List:  info,
	}, nil
}
