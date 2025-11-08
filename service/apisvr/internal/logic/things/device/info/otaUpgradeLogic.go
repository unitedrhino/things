package info

import (
	"context"

	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaUpgradeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 设备升级,获取升级包手动升级
func NewOtaUpgradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaUpgradeLogic {
	return &OtaUpgradeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OtaUpgradeLogic) OtaUpgrade(req *types.DeviceOtaUpgradeReq) (resp *types.DeviceOtaUpgradeResp, err error) {
	ret, err := l.svcCtx.DeviceM.DeviceOtaUpgrade(l.ctx, utils.Copy[dm.DeviceOtaUpgradeReq](req))
	return utils.Copy[types.DeviceOtaUpgradeResp](ret), err
}
