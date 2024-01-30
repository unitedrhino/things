package datamanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAreaApplyIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAreaApplyIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAreaApplyIndexLogic {
	return &UserAreaApplyIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserAreaApplyIndexLogic) UserAreaApplyIndex(in *sys.UserAreaApplyIndexReq) (*sys.UserAreaApplyIndexResp, error) {
	f := relationDB.UserAreaApplyFilter{AuthTypes: in.AuthTypes}
	total, err := relationDB.NewUserAreaApplyRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	list, err := relationDB.NewUserAreaApplyRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	return &sys.UserAreaApplyIndexResp{List: ToUserAreaApplyInfos(list), Total: total}, nil

}
