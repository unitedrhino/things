package auth

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthApiMultiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthApiMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthApiMultiUpdateLogic {
	return &AuthApiMultiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthApiMultiUpdateLogic) AuthApiMultiUpdate(req *types.AuthApiMultiUpdateReq) error {
	m := make([]*sys.AuthApiInfo, 0, len(req.List))
	for _, v := range req.List {
		m = append(m, &sys.AuthApiInfo{
			Route:  v.Route,
			Method: v.Method,
		})
	}

	resp, err := l.svcCtx.AuthRpc.AuthApiMultiUpdate(l.ctx, &sys.AuthApiMultiUpdateReq{
		RoleID: req.RoleID,
		List:   m,
	})

	if err != nil {
		l.Errorf("%s.rpc.AuthApiMultiUpdate req=%v err=%v", utils.FuncName(), req, err)
		return err
	}
	if resp == nil {
		l.Errorf("%s.rpc.AuthApiMultiUpdate return nil req=%v", utils.FuncName(), req)
		return errors.System.AddDetail("AuthApiMultiUpdate rpc return nil")
	}
	return nil
}
