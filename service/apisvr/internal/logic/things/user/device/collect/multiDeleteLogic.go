package collect

import (
	"context"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiDeleteLogic {
	return &MultiDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiDeleteLogic) MultiDelete(req *types.UserCollectDeviceSave) error {
	_, err := l.svcCtx.UserDevice.UserDeviceCollectMultiDelete(l.ctx, ToUserCollectDeviceSavePb(req))

	return err
}
