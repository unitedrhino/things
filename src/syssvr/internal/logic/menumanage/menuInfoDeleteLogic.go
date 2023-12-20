package menumanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	MiDB   *relationDB.MenuInfoRepo
	logx.Logger
}

func NewMenuInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuInfoDeleteLogic {
	return &MenuInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		MiDB:   relationDB.NewMenuInfoRepo(ctx),
	}
}

func (l *MenuInfoDeleteLogic) MenuInfoDelete(in *sys.WithID) (*sys.Response, error) {
	err := l.MiDB.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &sys.Response{}, nil
}
