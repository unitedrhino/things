package info

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindTokenCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建绑定token
func NewBindTokenCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindTokenCreateLogic {
	return &BindTokenCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindTokenCreateLogic) BindTokenCreate() (resp *types.DeviceBindTokenInfo, err error) {
	ret, err := l.svcCtx.DeviceM.DeviceBindTokenCreate(l.ctx, &dm.Empty{})
	return utils.Copy[types.DeviceBindTokenInfo](ret), err
}
