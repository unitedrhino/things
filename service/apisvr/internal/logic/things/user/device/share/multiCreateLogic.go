package share

import (
	"context"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 生成批量分享设备二维码
func NewMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiCreateLogic {
	return &MultiCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiCreateLogic) MultiCreate(req *types.UserMultiDevicesShareInfo) (resp *types.UserMultiDevicesShareKey, err error) {
	//将需要分享的设备记录缓存
	ret, err := l.svcCtx.UserDevice.UserMultiDevicesShareCreate(l.ctx, ToMuitlSharePb(req))
	if err != nil {
		return nil, err
	}
	return &types.UserMultiDevicesShareKey{ShareKey: ret.Key}, err
}
