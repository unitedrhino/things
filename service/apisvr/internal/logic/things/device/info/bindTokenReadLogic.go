package info

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindTokenReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 绑定token状态查询
func NewBindTokenReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindTokenReadLogic {
	return &BindTokenReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindTokenReadLogic) BindTokenRead(req *types.DeviceBindTokenReadReq) (resp *types.DeviceBindTokenInfo, err error) {
	ret, err := l.svcCtx.DeviceM.DeviceBindTokenRead(l.ctx, utils.Copy[dm.DeviceBindTokenReadReq](req))
	return utils.Copy[types.DeviceBindTokenInfo](ret), err
}
