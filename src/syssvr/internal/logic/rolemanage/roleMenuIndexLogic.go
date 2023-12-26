package rolemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleMenuIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	RiDB *relationDB.RoleInfoRepo
	MiDB *relationDB.MenuInfoRepo
}

func NewRoleMenuIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleMenuIndexLogic {
	return &RoleMenuIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		RiDB:   relationDB.NewRoleInfoRepo(ctx),
		MiDB:   relationDB.NewMenuInfoRepo(ctx),
	}
}

func (l *RoleMenuIndexLogic) RoleMenuIndex(in *sys.RoleMenuIndexReq) (*sys.RoleMenuIndexResp, error) {
	ms, err := relationDB.NewRoleMenuRepo(l.ctx).FindByFilter(l.ctx,
		relationDB.RoleMenuFilter{RoleIDs: []int64{in.Id}, AppCode: in.AppCode}, nil)
	if err != nil {
		return nil, err
	}
	if len(ms) == 0 { //没有菜单分配
		return &sys.RoleMenuIndexResp{}, nil
	}
	var menuIDs []int64
	if len(ms) != 0 {
		for _, v := range ms {
			menuIDs = append(menuIDs, v.MenuID)
		}
	}
	//menuInfos, err := l.MiDB.FindByFilter(l.ctx, relationDB.MenuInfoFilter{MenuIDs: menuIDs}, nil)
	//if err != nil {
	//	l.Errorf("MenuIndex find menu_info err,menuIds:%v,err:%v", menuIDs, err)
	//	return nil, err
	//}
	//for _, v := range menuInfos {
	//	info = append(info, logic.MenuInfoToPb(v))
	//}
	////查看缺失的父菜单Id
	//missingMenuInfos := logic.CheckMissingParentIdMenuIndex(l.ctx, menuInfos)
	//if len(missingMenuInfos) > 0 {
	//	for _, v := range missingMenuInfos {
	//		info = append(info, logic.MenuInfoToPb(v))
	//	}
	//}
	return &sys.RoleMenuIndexResp{MenuIDs: menuIDs}, nil
}
