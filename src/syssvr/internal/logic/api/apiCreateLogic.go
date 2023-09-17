package apilogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.ApiInfoRepo
}

func NewApiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiCreateLogic {
	return &ApiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewApiInfoRepo(ctx),
	}
}

func (l *ApiCreateLogic) ApiCreate(in *sys.ApiCreateReq) (*sys.Response, error) {
	err := l.AiDB.Insert(l.ctx, &relationDB.SysApiInfo{
		Route:        in.Route,
		Method:       in.Method,
		Name:         in.Name,
		BusinessType: in.BusinessType,
		Group:        in.Group,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &sys.Response{}, nil
}
