package apilogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.ApiInfoRepo
}

func NewApiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiUpdateLogic {
	return &ApiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewApiInfoRepo(ctx),
	}
}

func (l *ApiUpdateLogic) ApiUpdate(in *sys.ApiUpdateReq) (*sys.Response, error) {
	res, err := l.AiDB.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}

	err = l.AiDB.Update(l.ctx, &relationDB.SysApiInfo{
		ID:           in.Id,
		Route:        in.Route,
		Method:       in.Method,
		Name:         in.Name,
		Group:        in.Group,
		BusinessType: res.BusinessType,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &sys.Response{}, nil
}
