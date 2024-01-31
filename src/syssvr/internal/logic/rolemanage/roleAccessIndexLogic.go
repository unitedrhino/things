package rolemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleAccessIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleAccessIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAccessIndexLogic {
	return &RoleAccessIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleAccessIndexLogic) RoleAccessIndex(in *sys.RoleAccessIndexReq) (*sys.RoleAccessIndexResp, error) {
	ms, err := relationDB.NewRoleAccessRepo(l.ctx).FindByFilter(l.ctx,
		relationDB.RoleAccessFilter{RoleIDs: []int64{in.Id}}, nil)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, v := range ms {
		ids = append(ids, v.AccessCode)
	}
	return &sys.RoleAccessIndexResp{AccessCodes: ids}, nil
}
