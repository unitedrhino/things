package version

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.DeviceModuleVersionIndexReq) (resp *types.DeviceModuleVersionIndexResp, err error) {
	ret, err := l.svcCtx.DeviceM.DeviceModuleVersionIndex(l.ctx, utils.Copy[dm.DeviceModuleVersionIndexReq](req))

	return utils.Copy[types.DeviceModuleVersionIndexResp](ret), err
}
