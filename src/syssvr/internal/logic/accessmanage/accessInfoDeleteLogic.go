package accessmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccessInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAccessInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccessInfoDeleteLogic {
	return &AccessInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AccessInfoDeleteLogic) AccessInfoDelete(in *sys.WithID) (*sys.Response, error) {
	err := relationDB.NewAccessRepo(l.ctx).Delete(l.ctx, in.Id)
	return &sys.Response{}, err
}
