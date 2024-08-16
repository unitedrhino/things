package info

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiBindLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiBindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiBindLogic {
	return &MultiBindLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiBindLogic) MultiBind(req *types.DeviceInfoMultiBindReq) (resp *types.DeviceInfoMultiBindResp, err error) {
	ret, err := l.svcCtx.DeviceM.DeviceInfoMultiBind(l.ctx, utils.Copy[dm.DeviceInfoMultiBindReq](req))
	return utils.Copy[types.DeviceInfoMultiBindResp](ret), err
}
