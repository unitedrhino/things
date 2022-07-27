package device

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoDeleteLogic {
	return &InfoDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoDeleteLogic) InfoDelete(req *types.DeviceInfoDeleteReq) error {
	_, err := l.svcCtx.DmRpc.DeviceInfoDelete(l.ctx, &dm.DeviceInfoDeleteReq{
		ProductID:  req.ProductID,
		DeviceName: req.DeviceName,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.ManageDevice|req=%v|err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
