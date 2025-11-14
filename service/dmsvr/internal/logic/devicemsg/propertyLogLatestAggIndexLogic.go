package devicemsglogic

import (
	"context"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyLogLatestAggIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPropertyLogLatestAggIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyLogLatestAggIndexLogic {
	return &PropertyLogLatestAggIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PropertyLogLatestAggIndexLogic) PropertyLogLatestAggIndex(in *dm.PropertyLatestAggIndexReq) (*dm.PropertyLatestAggIndexResp, error) {
	// todo: add your logic here and delete this line

	return &dm.PropertyLatestAggIndexResp{}, nil
}
