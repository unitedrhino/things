package user

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/spf13/cast"

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

func (l *DeleteLogic) Delete(req *types.UserDeleteReq) error {
	_, err := l.svcCtx.UserRpc.Delete(l.ctx, &sys.UserDeleteReq{
		Uid: cast.ToInt64(req.Uid)})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.InfoDelete req=%v err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
