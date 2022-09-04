package menu

import (
	"context"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMenuIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuIndexLogic {
	return &MenuIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenuIndexLogic) MenuIndex(req *types.MenuIndexReq) (resp *types.MenuIndexResp, err error) {
	var page sys.PageInfo
	copier.Copy(&page, req.Page)
	info, err := l.svcCtx.MenuRpc.MenuIndex(l.ctx, &sys.MenuIndexReq{
		Page: &page,
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}

	var menuInfo []*types.MenuIndexData
	var total int64
	total = info.Total

	menuInfo = make([]*types.MenuIndexData, 0, len(menuInfo))
	for _, i := range info.List {
		menuInfo = append(menuInfo, &types.MenuIndexData{
			ID:         i.Id,
			Name:       i.Name,
			ParentID:   i.ParentID,
			Type:       i.Type,
			Path:       i.Path,
			Component:  i.Component,
			Icon:       i.Icon,
			Redirect:   i.Redirect,
			CreateTime: i.CreateTime,
			Order:      i.Order,
		})
	}

	return &types.MenuIndexResp{menuInfo, total}, nil
}
