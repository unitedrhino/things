package share

import (
	"context"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiDeleteTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除批量分享 Token
func NewMultiDeleteTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiDeleteTokenLogic {
	return &MultiDeleteTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiDeleteTokenLogic) MultiDeleteToken(req *types.UserDeviceShareMultiDeleteTokenReq) error {
	_, err := l.svcCtx.UserDevice.UserDeviceShareMultiDeleteToken(l.ctx, &dm.UserDeviceShareMultiDeleteTokenReq{
		ShareToken: req.ShareToken,
	})
	return err
}
