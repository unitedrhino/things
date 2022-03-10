package user

import (
	"context"
	"github.com/go-things/things/src/usersvr/user"
	"github.com/jinzhu/copier"

	"github.com/go-things/things/src/webapi/internal/svc"
	"github.com/go-things/things/src/webapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCoreListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCoreListLogic(ctx context.Context, svcCtx *svc.ServiceContext) UserCoreListLogic {
	return UserCoreListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCoreListLogic) UserCoreList(req types.GetUserCoreListReq) (*types.GetUserCoreListResp, error) {
	l.Infof("UserCoreList|req=%+v", req)
	var page user.PageInfo
	copier.Copy(&page, req.Page)
	info, err := l.svcCtx.UserRpc.GetUserCoreList(l.ctx, &user.GetUserCoreListReq{
		Page: &page,
	})
	if err != nil {
		return nil, err
	}
	resp := types.GetUserCoreListResp{
		Total: info.Total,
	}
	resp.Info = make([]*types.UserCore, 0, len(resp.Info))
	for _, i := range info.Info {
		resp.Info = append(resp.Info, types.UserCoreToApi(i))
	}
	return &resp, nil
}
