package user

import (
	"context"
	"github.com/i-Things/things/src/apisvr/internal/domain/userHeader"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
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
	menuInfo := make([]*types.MenuData, 0)
	info, err := l.svcCtx.MenuRpc.MenuIndex(l.ctx, &sys.MenuIndexReq{
		Role: userHeader.GetUserCtx(l.ctx).Role,
	})
	if err != nil {
		return nil, err
	}

	for _, me := range info.List {
		menuInfo = append(menuInfo, &types.MenuData{
			ID:         me.Id,
			Name:       me.Name,
			ParentID:   me.ParentID,
			Type:       me.Type,
			Path:       me.Path,
			Component:  me.Component,
			Icon:       me.Icon,
			Redirect:   me.Redirect,
			CreateTime: me.CreateTime,
			Order:      me.Order,
			HideInMenu: me.HideInMenu,
		})
	}

	return &types.UserResourceReadResp{
		Menu: menuInfo,
	}, nil

}
