package share

import (
	"context"
	"gitee.com/i-Things/things/service/apisvr/internal/logic"
	"gitee.com/i-Things/things/service/apisvr/internal/logic/things"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"gitee.com/i-Things/things/service/apisvr/internal/types"

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

func (l *IndexLogic) Index(req *types.UserDeviceShareIndexReq) (resp *types.UserDeviceShareIndexResp, err error) {
	ret, err := l.svcCtx.UserDevice.UserDeviceShareIndex(l.ctx, &dm.UserDeviceShareIndexReq{
		Page:   logic.ToDmPageRpc(req.Page),
		Device: things.ToDmDeviceCorePb(req.Device),
	})
	if err != nil {
		return nil, err
	}

	return &types.UserDeviceShareIndexResp{
		List:  ToSharesTypes(ret.List),
		Total: ret.Total,
	}, nil
}
