package dataUpdateEvent

import (
	"context"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/ddsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DataUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDataUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataUpdateLogic {
	return &DataUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (d *DataUpdateLogic) ProductCustomUpdate(info *events.DataUpdateInfo) error {
	d.Infof("%s DataUpdateInfo:%v", utils.FuncName(), info)
	return d.svcCtx.Script.ClearCache(d.ctx, info.ProductID)
}
