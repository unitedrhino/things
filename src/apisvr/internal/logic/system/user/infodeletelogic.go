package user

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/users"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/middleware"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/usersvr/user"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

type InfoDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoDeleteLogic {
	return &InfoDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoDeleteLogic) InfoDelete(req *types.UserInfoDeleteReq) error {
	middleware.NewRecordMiddleware()
	//从context中获取Uid再传入l.svcCtx.UserRpc.InfoDelete
	_, err := l.svcCtx.UserRpc.InfoDelete(l.ctx, &user.UserInfoDeleteReq{Uid: cast.ToString(users.GetClaimsFromToken(l.ctx, types.USER_UID, "").Uid)})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.InfoDelete|req=%v|err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
