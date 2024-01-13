package usermanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UiDB *relationDB.UserInfoRepo
}

func NewUserInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoReadLogic {
	return &UserInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
	}
}

func (l *UserInfoReadLogic) UserInfoRead(in *sys.UserInfoReadReq) (*sys.UserInfo, error) {
	if err := ctxs.IsRoot(l.ctx); err == nil {
		ctxs.GetUserCtx(l.ctx).AllTenant = true
		defer func() {
			ctxs.GetUserCtx(l.ctx).AllTenant = false
		}()
	}
	ui, err := l.UiDB.FindOne(l.ctx, in.UserID)
	if err != nil {
		l.Logger.Error("UserInfoModel.FindOne err:%v", err)
		return nil, err
	}

	return UserInfoToPb(l.ctx, ui, l.svcCtx), nil
}
