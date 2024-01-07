package rolemanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/zeromicro/go-zero/core/logx"
)

type RoleApiIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleApiIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleApiIndexLogic {
	return &RoleApiIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleApiIndexLogic) RoleApiIndex(in *sys.RoleApiIndexReq) (*sys.RoleApiIndexResp, error) {
	uc := ctxs.GetUserCtx(l.ctx)
	if uc.IsAdmin { //超级管理员默认全部勾选
		ms, err := relationDB.NewTenantAppApiRepo(l.ctx).FindByFilter(l.ctx, relationDB.TenantAppApiFilter{AppCode: in.AppCode, ModuleCode: in.ModuleCode}, nil)
		if err != nil {
			return nil, err
		}
		var ids []int64
		for _, v := range ms {
			ids = append(ids, v.ID)
		}
		return &sys.RoleApiIndexResp{ApiIDs: ids}, nil
	}
	ms, err := relationDB.NewRoleApiRepo(l.ctx).FindByFilter(l.ctx,
		relationDB.RoleApiFilter{RoleIDs: []int64{in.Id}, AppCode: in.AppCode, ModuleCode: in.ModuleCode}, nil)
	if err != nil {
		return nil, err
	}
	var ids []int64
	for _, v := range ms {
		ids = append(ids, v.ApiID)
	}
	return &sys.RoleApiIndexResp{ApiIDs: ids}, nil
}
