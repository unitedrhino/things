package stream

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/vidsvr/client/vidmgrinfomanage"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

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

func (l *DeleteLogic) Delete(req *types.VidmgrStreamDeleteReq) error {
	// todo: add your logic here and delete this line
	_, err := l.svcCtx.VidmgrS.VidmgrStreamDelete(l.ctx, &vidmgrinfomanage.VidmgrStreamDeleteReq{
		StreamID: req.StreamID,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.VidmgrStreamDelete req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
