package user

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic/system"
	app "github.com/i-Things/things/src/apisvr/internal/logic/system/app/info"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/role"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"golang.org/x/sync/errgroup"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResourceReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResourceReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceReadLogic {
	return &ResourceReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceReadLogic) ResourceRead() (resp *types.UserResourceReadResp, err error) {
	var (
		menuInfo []*types.MenuInfo
		userInfo *types.UserInfo
		projects []*types.ProjectInfo
		appInfo  *types.AppInfo
		roles    []*types.RoleInfo
		wait     errgroup.Group
		uc       = ctxs.GetUserCtx(l.ctx)
	)
	wait.Go(func() error {
		defer utils.Recover(l.ctx)
		info, err := l.svcCtx.RoleRpc.RoleMenuIndex(l.ctx, &sys.RoleMenuIndexReq{
			Id:      uc.RoleID,
			AppCode: uc.AppCode,
		})
		if err != nil {
			return err
		}
		for _, me := range info.List {
			menuInfo = append(menuInfo, system.ToMenuInfoApi(me))
		}
		return nil
	})
	wait.Go(func() error {
		defer utils.Recover(l.ctx)
		info, err := l.svcCtx.UserRpc.UserRoleIndex(l.ctx, &sys.UserRoleIndexReq{
			UserID: uc.UserID,
		})
		if err != nil {
			return err
		}
		roles = role.ToRoleInfosTypes(info.List)
		return nil
	})
	wait.Go(func() error {
		ui, err := l.svcCtx.UserRpc.UserInfoRead(l.ctx, &sys.UserInfoReadReq{UserID: ctxs.GetUserCtx(l.ctx).UserID})
		if err != nil {
			return err
		}
		userInfo = UserInfoToApi(ui)
		return nil
	})
	wait.Go(func() error {
		ret, err := l.svcCtx.AppRpc.AppInfoRead(l.ctx, &sys.WithIDCode{Code: ctxs.GetUserCtx(l.ctx).AppCode})
		if err != nil {
			return err
		}
		appInfo = app.ToAppInfoTypes(ret)
		return nil
	})
	wait.Go(func() error {
		pis, err := l.svcCtx.ProjectM.ProjectInfoIndex(l.ctx, &sys.ProjectInfoIndexReq{})
		if err != nil {
			return err
		}
		for _, pb := range pis.List {
			projects = append(projects, system.ProjectInfoToApi(pb))
		}
		return nil
	})
	err = wait.Wait()
	if err != nil {
		return nil, err
	}
	return &types.UserResourceReadResp{
		Menus:    menuInfo,
		Info:     userInfo,
		App:      appInfo,
		Roles:    roles,
		Projects: projects,
	}, nil
}
