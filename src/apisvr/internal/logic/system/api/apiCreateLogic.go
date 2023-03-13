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

type ApiCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiCreateLogic {
	return &ApiCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiCreateLogic) ApiCreate(req *types.ApiCreateReq) error {
	resp, err := l.svcCtx.ApiRpc.ApiCreate(l.ctx, &sys.ApiCreateReq{
		Route:        req.Route,
		Method:       req.Method,
		Group:        req.Group,
		Name:         req.Name,
		BusinessType: req.BusinessType,
	})
	if err != nil {
		err := errors.Fmt(err)
		l.Errorf("%s.rpc.ApiCreate req=%v err=%+v", utils.FuncName(), req, err)
		return err
	}
	if resp == nil {
		l.Errorf("%s rpc.ApiCreate return nil req=%+v", utils.FuncName(), req)
		return errors.System.AddDetail("ApiCreate rpc return nil")
	}
	return nil
}
