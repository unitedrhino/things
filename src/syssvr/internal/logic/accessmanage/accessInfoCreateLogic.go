package accessmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccessInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAccessInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccessInfoCreateLogic {
	return &AccessInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AccessInfoCreateLogic) AccessInfoCreate(in *sys.AccessInfo) (*sys.WithID, error) {
	po := ToAccessPo(in)
	po.ID = 0
	err := relationDB.NewAccessRepo(l.ctx).Insert(l.ctx, po)
	if err != nil {
		return nil, err
	}
	return &sys.WithID{Id: po.ID}, nil
}
