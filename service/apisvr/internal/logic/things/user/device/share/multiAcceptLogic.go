package share

import (
	"context"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiAcceptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 接受批量分享设备
func NewMultiAcceptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiAcceptLogic {
	return &MultiAcceptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiAcceptLogic) MultiAccept(req *types.UserDeviceShareMultiAcceptInfo) error {
	uc := ctxs.GetUserCtx(l.ctx)
	dmreq := dm.UserDeviceShareMultiAcceptReq{
		ShareToken:        req.ShareToken,
		SharedUserAccount: uc.UserName,
		SharedUserID:      uc.UserID,
		Devices:           ToSharesDevices(req.Devices),
	}
	_, err := l.svcCtx.UserDevice.UserDeivceShareMultiAccept(l.ctx, &dmreq)
	if err != nil {
		return err
	}
	return nil
}
