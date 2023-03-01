package apilogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiDeleteLogic {
	return &ApiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiDeleteLogic) ApiDelete(in *sys.ApiDeleteReq) (*sys.Response, error) {
	err := l.svcCtx.ApiModel.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &sys.Response{}, nil
}
