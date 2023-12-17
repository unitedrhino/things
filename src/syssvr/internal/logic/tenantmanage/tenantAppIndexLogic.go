package tenantmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	appmanagelogic "github.com/i-Things/things/src/syssvr/internal/logic/appmanage"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppIndexLogic {
	return &TenantAppIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppIndexLogic) TenantAppIndex(in *sys.TenantAppIndexReq) (*sys.TenantAppIndexResp, error) {
	uc := ctxs.GetUserCtx(l.ctx)
	if uc.TenantCode != def.TenantCodeDefault {
		in.Code = uc.TenantCode
	} else {
		uc.AllData = true
		defer func() { uc.AllData = false }()
	}
	f := relationDB.TenantAppFilter{TenantCode: in.Code}
	list, err := relationDB.NewTenantAppRepo(l.ctx).FindByFilter(l.ctx, f, nil)
	if err != nil {
		return nil, err
	}
	appCodes := make([]string, 0)
	for _, v := range list {
		appCodes = append(appCodes, v.AppCode)
	}
	apps, err := relationDB.NewAppInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.AppInfoFilter{Codes: appCodes}, nil)
	if err != nil {
		return nil, err
	}
	return &sys.TenantAppIndexResp{List: appmanagelogic.ToAppInfosPb(apps), Total: int64(len(list))}, nil
}
