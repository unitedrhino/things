package dataUpdateEvent

import (
	"context"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DataUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataUpdateLogic {
	return &DataUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (d *DataUpdateLogic) ProductSchemaUpdate(info *events.DataUpdateInfo) error {
	d.Infof("%s DataUpdateInfo:%v", utils.FuncName(), info)
	return d.svcCtx.SchemaRepo.ClearCache(d.ctx, info.ProductID)
}
