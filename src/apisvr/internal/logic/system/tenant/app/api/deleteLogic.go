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

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.WithAppCodeID) error {
	_, err := l.svcCtx.TenantRpc.TenantAppApiDelete(l.ctx, &sys.WithAppCodeID{
		AppCode: req.AppCode,
		Code:    req.Code,
		Id:      req.ID,
	})
	if err != nil {
		err := errors.Fmt(err)
		l.Errorf("%s.rpc.ApiDelete req=%v err=%+v", utils.FuncName(), req, err)
		return err
	}
	return nil
}