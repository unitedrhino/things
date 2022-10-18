package menulogic

import (
	"context"
	"database/sql"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/mysql"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuCreateLogic {
	return &MenuCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MenuCreateLogic) MenuCreate(in *sys.MenuCreateReq) (*sys.Response, error) {
	_, err := l.svcCtx.MenuInfoModle.Insert(l.ctx, &mysql.MenuInfo{
		ParentID:      sql.NullInt64{Int64: in.ParentID, Valid: true},
		Type:          sql.NullInt64{Int64: in.Type, Valid: true},
		Order:         sql.NullInt64{Int64: in.Order, Valid: true},
		Name:          in.Name,
		Path:          in.Path,
		Component:     in.Component,
		Icon:          in.Icon,
		Redirect:      in.Redirect,
		BackgroundUrl: "",
		HideInMenu:    in.HideInMenu,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &sys.Response{}, nil
}
