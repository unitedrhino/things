package menulogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuIndexLogic {
	return &MenuIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MenuIndexLogic) MenuIndex(in *sys.MenuIndexReq) (*sys.MenuIndexResp, error) {
	mes, total, err := l.svcCtx.MenuModel.Index(in)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	info := make([]*sys.MenuIndexData, 0, len(mes))
	for _, me := range mes {
		info = append(info, &sys.MenuIndexData{
			Id:         me.Id,
			Name:       me.Name,
			ParentID:   me.ParentID,
			Type:       me.Type,
			Path:       me.Path,
			Component:  me.Component,
			Icon:       me.Icon,
			Redirect:   me.Redirect,
			CreateTime: me.CreatedTime.Unix(),
			Order:      me.Order,
		})
	}
	return &sys.MenuIndexResp{
		List:  info,
		Total: total,
	}, nil

	return &sys.MenuIndexResp{}, nil
}
