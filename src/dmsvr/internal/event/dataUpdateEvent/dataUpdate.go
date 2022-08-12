package dataUpdateEvent

import (
	"context"
	"github.com/i-Things/things/shared/domain/schema"
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

func (d DataUpdateLogic) SchemaClearCache(info *schema.SchemaInfo) error {
	d.Infof("DataUpdateLogic|SchemaClearCache|productID:%v", info.ProductID)
	return d.svcCtx.SchemaRepo.ClearCache(d.ctx, info.ProductID)
}
