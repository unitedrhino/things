package dataUpdateEvent

import (
	"context"
	"gitee.com/i-Things/core/shared/events"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/rulesvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/rulesvr/internal/svc"
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

func (d *DataUpdateLogic) ProductSchemaUpdate(info *events.DeviceUpdateInfo) error {
	d.Infof("%s DeviceUpdateInfo:%v", utils.FuncName(), info)
	return d.svcCtx.SchemaRepo.ClearCache(d.ctx, info.ProductID)
}

func (d *DataUpdateLogic) SceneInfoDelete(info *events.ChangeInfo) error {
	d.Infof("%s DeviceUpdateInfo:%v", utils.FuncName(), info)
	return d.svcCtx.SceneTimerControl.Delete(info.ID)
}
func (d *DataUpdateLogic) SceneInfoUpdate(info *events.ChangeInfo) error {
	d.Infof("%s DeviceUpdateInfo:%v", utils.FuncName(), info)
	do, err := relationDB.NewSceneInfoRepo(d.ctx).FindOne(d.ctx, info.ID)
	if err != nil {
		return err
	}
	return d.svcCtx.SceneTimerControl.Update(do)
}
