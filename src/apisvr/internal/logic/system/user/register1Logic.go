package user

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Register1Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegister1Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Register1Logic {
	return &Register1Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Register1Logic) Register1(req *types.UserRegister1Req) (*types.UserRegister1Resp, error) {
	resp, err := l.svcCtx.UserRpc.UserRegister1(l.ctx, &sys.UserRegister1Req{
		RegType: req.RegType,
		Note:    req.Note,
		Code:    req.Code,
		CodeID:  req.CodeID,
	})
	if err != nil {
		l.Errorf("%s.rpc.Register1 req=%v err=%v ", utils.FuncName(), req, err)
		return &types.UserRegister1Resp{}, err
	}
	return &types.UserRegister1Resp{Token: resp.Token}, nil
}
