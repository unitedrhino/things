package usermanagelogic

import (
	"context"
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
	ui, err := l.UiDB.FindOne(l.ctx, in.UserID)
	if err != nil {
		l.Logger.Error("UserInfoModel.FindOne err , sql:%s", l.svcCtx)
		return nil, err
	}

	return UserInfoToPb(ui), nil
}
