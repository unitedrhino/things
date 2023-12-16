package apimanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.ApiInfoRepo
}

func NewApiInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiInfoDeleteLogic {
	return &ApiInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewApiInfoRepo(ctx),
	}
}

func (l *ApiInfoDeleteLogic) ApiInfoDelete(in *sys.ReqWithID) (*sys.Response, error) {
	err := l.AiDB.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &sys.Response{}, nil
}
