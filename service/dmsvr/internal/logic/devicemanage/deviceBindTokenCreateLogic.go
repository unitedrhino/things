package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceBind"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceBindTokenCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceBindTokenCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceBindTokenCreateLogic {
	return &DeviceBindTokenCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceBindTokenCreateLogic) DeviceBindTokenCreate(in *dm.Empty) (*dm.DeviceBindTokenInfo, error) {
	shareToken := devices.GenMsgToken(l.ctx, l.svcCtx.NodeID)
	uc := ctxs.GetUserCtx(l.ctx)
	bt := deviceBind.TokenInfo{
		Token:  shareToken,
		UserID: uc.UserID,
		Status: deviceBind.StatusInit,
	}
	err := l.svcCtx.DeviceBindToken.SetData(l.ctx, shareToken, &bt)
	return &dm.DeviceBindTokenInfo{
		Token:  shareToken,
		Status: deviceBind.StatusInit,
	}, err
}
