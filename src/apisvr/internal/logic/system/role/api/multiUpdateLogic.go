package api

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiUpdateLogic {
	return &MultiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiUpdateLogic) MultiUpdate(req *types.RoleApiMultiUpdateReq) error {
	m := make([]*sys.AuthApiInfo, 0, len(req.Apis))
	for _, v := range req.Apis {
		m = append(m, &sys.AuthApiInfo{
			Route:  v.Route,
			Method: v.Method,
		})
	}

	resp, err := l.svcCtx.RoleRpc.RoleApiMultiUpdate(l.ctx, &sys.RoleApiMultiUpdateReq{
		Id:      req.ID,
		AppCode: req.AppCode,
		List:    m,
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
