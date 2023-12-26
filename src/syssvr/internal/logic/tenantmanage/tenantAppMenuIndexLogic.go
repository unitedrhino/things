package tenantmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppMenuIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppMenuIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppMenuIndexLogic {
	return &TenantAppMenuIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppMenuIndexLogic) TenantAppMenuIndex(in *sys.TenantAppMenuIndexReq) (*sys.TenantAppMenuIndexResp, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	f := relationDB.TenantAppMenuFilter{
		ModuleCode: in.ModuleCode,
		TenantCode: in.Code,
		AppCode:    in.AppCode,
	}
	resp, err := relationDB.NewTenantAppMenuRepo(l.ctx).FindByFilter(l.ctx, f, nil)
	if err != nil {
		return nil, err
	}
	//total, err := relationDB.NewMenuInfoRepo(l.ctx).CountByFilter(l.ctx, f)
	//if err != nil {
	//	return nil, err
	//}
	info := make([]*sys.TenantAppMenu, 0, len(resp))
	if !in.IsRetTree {
		for _, v := range resp {
			i := logic.ToTenantAppMenuInfoPb(v)
			info = append(info, i)
		}
		return &sys.TenantAppMenuIndexResp{List: info}, nil
	}

	var (
		pidMap = make(map[int64][]*sys.TenantAppMenu, len(resp))
		idMap  = make(map[int64]*sys.TenantAppMenu, len(resp))
	)
	for _, v := range resp {
		i := logic.ToTenantAppMenuInfoPb(v)
		idMap[i.Info.Id] = i
		if i.Info.ParentID == 1 { //根节点
			info = append(info, i)
			continue
		}
		pidMap[i.Info.ParentID] = append(pidMap[i.Info.ParentID], i)
	}
	fillChildren(info, pidMap)
	return &sys.TenantAppMenuIndexResp{List: info}, nil
}

func fillChildren(in []*sys.TenantAppMenu, pidMap map[int64][]*sys.TenantAppMenu) {
	for _, v := range in {
		cs := pidMap[v.Info.Id]
		if cs != nil {
			v.Children = cs
			fillChildren(cs, pidMap)
		}
	}
}
