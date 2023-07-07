package userlogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UiDB *relationDB.UserInfoRepo
}

func NewUserReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
	}
}

func (l *ReadLogic) UserRead(in *sys.UserReadReq) (*sys.UserInfo, error) {
	ui, err := l.UiDB.FindOne(l.ctx, in.UserID)
	if err != nil {
		l.Logger.Error("UserInfoModel.FindOne err , sql:%s", l.svcCtx)
		return nil, err
	}

	return UserInfoToPb(ui), nil
}
