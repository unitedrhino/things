package usermanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAreaMultiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAreaMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAreaMultiDeleteLogic {
	return &UserAreaMultiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserAreaMultiDeleteLogic) UserAreaMultiDelete(in *sys.UserAreaMultiDeleteReq) (*sys.Response, error) {
	err := relationDB.NewUserAreaRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.UserAreaFilter{AreaIDs: in.AreaIDs})
	return &sys.Response{}, err
}
