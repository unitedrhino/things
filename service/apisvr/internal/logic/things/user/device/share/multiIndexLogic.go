package share

import (
	"context"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取批量分享的设备列表
func NewMultiIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiIndexLogic {
	return &MultiIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiIndexLogic) MultiIndex(req *types.UserMultiDevicesShareToken) (resp *types.UserMultiDevicesShareInfo, err error) {

	ret, err := l.svcCtx.UserDevice.UserMultiDeivcesShareIndex(l.ctx, &dm.UserMultiDevicesShareKeyword{ShareToken: req.ShareToken})
	if err != nil {
		return nil, err
	}
	return ToMultiShareTypes(ret), err
}
