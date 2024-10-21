package userdevicelogic

import (
	"context"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeivceShareMultiIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeivceShareMultiIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeivceShareMultiIndexLogic {
	return &UserDeivceShareMultiIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 扫码后获取设备列表
func (l *UserDeivceShareMultiIndexLogic) UserDeivceShareMultiIndex(in *dm.UserDeviceShareMultiToken) (*dm.UserDeviceShareMultiInfo, error) {
	resp, err := l.svcCtx.UserMultiDeviceShare.GetData(l.ctx, in.ShareToken)
	return resp, err
}
