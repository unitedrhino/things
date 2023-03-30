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

type ApiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiUpdateLogic {
	return &ApiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiUpdateLogic) ApiUpdate(req *types.ApiUpdateReq) error {
	resp, err := l.svcCtx.ApiRpc.ApiUpdate(l.ctx, &sys.ApiUpdateReq{
		Id:     req.ID,
		Route:  req.Route,
		Method: req.Method,
		Group:  req.Group,
		Name:   req.Name,
	})
	if err != nil {
		err := errors.Fmt(err)
		l.Errorf("%s.rpc.ApiUpdate req=%v err=%+v", utils.FuncName(), req, err)
		return err
	}
	if resp == nil {
		l.Errorf("%s.rpc.ApiUpdate return nil req=%+v", utils.FuncName(), req)
		return errors.System.AddDetail("ApiUpdate rpc return nil")
	}

	return nil
}
