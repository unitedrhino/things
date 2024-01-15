package self

import (
	"context"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ForgetPwdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewForgetPwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForgetPwdLogic {
	return &ForgetPwdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ForgetPwdLogic) ForgetPwd(req *types.UserForgetPwdReq) error {
	_, err := l.svcCtx.UserRpc.UserForgetPwd(l.ctx, &sys.UserForgetPwdReq{
		Account:  req.Account,
		Type:     req.Type,
		Password: req.Password,
		Code:     req.Code,
		CodeID:   req.CodeID,
	})

	return err
}
