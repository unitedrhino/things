package accessmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiInfoUpdateLogic {
	return &ApiInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiInfoUpdateLogic) ApiInfoUpdate(in *sys.ApiInfo) (*sys.Response, error) {
	old, err := relationDB.NewApiInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	old.AccessCode = in.AccessCode
	old.Method = in.Method
	old.Route = in.Route
	old.Name = in.Name
	old.BusinessType = in.BusinessType
	old.Desc = in.Desc
	old.IsAuthTenant = in.IsAuthTenant
	err = relationDB.NewApiInfoRepo(l.ctx).Update(l.ctx, old)
	return &sys.Response{}, err
}
