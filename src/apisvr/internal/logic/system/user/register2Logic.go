package user

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Register2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegister2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Register2Logic {
	return &Register2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Register2Logic) Register2(req *types.UserRegister2Req) error {
	if req.UserInfo.Sex != 1 && req.UserInfo.Sex != 2 {
		req.UserInfo.Sex = 1
	}
	_, err := l.svcCtx.UserRpc.UserRegister2(l.ctx, &sys.UserRegister2Req{
		Token: req.Token,
		RegIP: "",
		Info:  UserInfoToRpc(&req.UserInfo),
	})
	if err != nil {
		l.Errorf("%s.rpc.Register1 req=%v err=%v ", utils.FuncName(), req, err)
		return err
	}
	return nil
}
