package apimanagelogic

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiInfoCreateLogic {
	return &ApiInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiInfoCreateLogic) ApiInfoCreate(in *sys.ApiInfo) (*sys.WithID, error) {
	// todo: add your logic here and delete this line

	return &sys.WithID{}, nil
}
