package dataUpdateEvent

import (
	"context"
	"github.com/i-Things/things/src/dmsvr/internal/domain/thing"
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

func (d DataUpdateLogic) TempModelClearCache(info *thing.TemplateInfo) error {
	d.Infof("DataUpdateLogic|TempModelClearCache|productID:%v", info.ProductID)
	return d.svcCtx.TemplateRepo.ClearCache(d.ctx, info.ProductID)
}
