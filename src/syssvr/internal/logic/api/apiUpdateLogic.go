package apilogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/mysql"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiUpdateLogic {
	return &ApiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiUpdateLogic) ApiUpdate(in *sys.ApiUpdateReq) (*sys.Response, error) {
	res, err := l.svcCtx.ApiModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}

	err = l.svcCtx.ApiModel.Update(l.ctx, &mysql.SysApi{
		Id:           in.Id,
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
