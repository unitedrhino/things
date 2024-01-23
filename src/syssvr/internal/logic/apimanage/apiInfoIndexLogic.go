package apimanagelogic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &sys.ApiInfoIndexResp{}, nil
}
