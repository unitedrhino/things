package usermanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/spf13/cast"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UiDB *relationDB.UserInfoRepo
}

func NewUserInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoDeleteLogic {
	return &UserInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
	}
}

func (l *UserInfoDeleteLogic) UserInfoDelete(in *sys.UserInfoDeleteReq) (*sys.Response, error) {
	err := l.UiDB.Delete(l.ctx, cast.ToInt64(in.UserID))
	if err != nil {
		l.Errorf("%s.Delete uid=%d err=%+v", utils.FuncName(), in.UserID, err)
		return nil, err
	}

	l.Infof("%s.delete uid=%v", utils.FuncName(), in.UserID)

	return &sys.Response{}, nil
}
