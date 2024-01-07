package rolemanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleAppIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleAppIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAppIndexLogic {
	return &RoleAppIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleAppIndexLogic) RoleAppIndex(in *sys.RoleAppIndexReq) (*sys.RoleAppIndexResp, error) {
	var appCodes []string
	uc := ctxs.GetUserCtx(l.ctx)
	if uc.IsAdmin { //超级管理员默认全部勾选
		ret, err := relationDB.NewTenantAppRepo(l.ctx).FindByFilter(l.ctx, relationDB.TenantAppFilter{}, nil)
		if err != nil {
			return nil, err
		}
		for _, v := range ret {
			appCodes = append(appCodes, v.AppCode)
		}
		return &sys.RoleAppIndexResp{AppCodes: appCodes, Total: int64(len(appCodes))}, nil
	}
	ra, err := relationDB.NewRoleAppRepo(l.ctx).FindByFilter(l.ctx, relationDB.RoleAppFilter{RoleID: in.Id}, nil)
	if err != nil {
		return nil, err
	}
	for _, v := range ra {
		appCodes = append(appCodes, v.AppCode)
	}
	return &sys.RoleAppIndexResp{AppCodes: appCodes, Total: int64(len(appCodes))}, nil
}
