package userdevicelogic

import (
	"context"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserMultiDeivcesShareIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserMultiDeivcesShareIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserMultiDeivcesShareIndexLogic {
	return &UserMultiDeivcesShareIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserMultiDeivcesShareIndexLogic) UserMultiDeivcesShareIndex(in *dm.UserMultiDevicesShareKeyword) (resp *dm.UserMultiDevicesShareInfo, err error) {
	resp, err = l.svcCtx.UserMultiDeviceShare.GetData(l.ctx, in.ShareToken)
	return resp, err
}
