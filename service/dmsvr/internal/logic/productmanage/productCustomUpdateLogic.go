package productmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/share/events"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductCustomUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PcDB *relationDB.ProductCustomRepo
}

func NewProductCustomUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductCustomUpdateLogic {
	return &ProductCustomUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PcDB:   relationDB.NewProductCustomRepo(ctx),
	}
}

func (l *ProductCustomUpdateLogic) ProductCustomUpdate(in *dm.ProductCustom) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	pi, err := l.PcDB.FindOneByProductID(l.ctx, in.ProductID)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			if in.ScriptLang == 0 {
				in.ScriptLang = 1
			}
			err = l.PcDB.Insert(l.ctx, &relationDB.DmProductCustom{
				ProductID:       in.ProductID,
				ScriptLang:      in.ScriptLang,
				CustomTopics:    logic.ToCustomTopicsDo(in.CustomTopics),
				TransformScript: in.TransformScript.GetValue(),
			})
			if err != nil {
				return nil, err
			}
			return &dm.Empty{}, nil
		}
		return nil, err
	}
	if in.TransformScript != nil {
		pi.TransformScript = in.TransformScript.GetValue()
	}
	if in.ScriptLang != 0 {
		pi.ScriptLang = in.ScriptLang
	}
	if in.CustomTopics != nil {
		pi.CustomTopics = logic.ToCustomTopicsDo(in.CustomTopics)
	}
	err = l.PcDB.Update(l.ctx, pi)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	err = l.svcCtx.FastEvent.Publish(l.ctx, eventBus.DmProductCustomUpdate, &events.DeviceUpdateInfo{ProductID: in.ProductID})
	if err != nil {
		return nil, err
	}
	return &dm.Empty{}, nil
}
