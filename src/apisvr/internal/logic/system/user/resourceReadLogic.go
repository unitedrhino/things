package user

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic/system"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/errgroup"
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
		menuInfo []*types.MenuData
		userInfo *types.UserInfo
		wait     errgroup.Group
	)
	wait.Go(func() error {
		defer utils.Recover(l.ctx)
		info, err := l.svcCtx.MenuRpc.MenuIndex(l.ctx, &sys.MenuIndexReq{
			Role: ctxs.GetUserCtx(l.ctx).Role,
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
		ui, err := l.svcCtx.UserRpc.UserRead(l.ctx, &sys.UserReadReq{UserID: ctxs.GetUserCtx(l.ctx).UserID})
		if err != nil {
			return err
		}
		userInfo = UserInfoToApi(ui)
		return nil
	})
	err = wait.Wait()
	if err != nil {
		return nil, err
	}
	return &types.UserResourceReadResp{
		Menu: menuInfo,
		Info: userInfo,
	}, nil
}
