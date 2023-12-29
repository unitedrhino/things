package apimanagelogic

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiInfoDeleteLogic {
	return &ApiInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiInfoDeleteLogic) ApiInfoDelete(in *sys.WithID) (*sys.Response, error) {
	// todo: add your logic here and delete this line

	return &sys.Response{}, nil
}
