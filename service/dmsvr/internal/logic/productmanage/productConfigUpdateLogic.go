package productmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/topics"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductConfigUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductConfigUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductConfigUpdateLogic {
	return &ProductConfigUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新产品配置
func (l *ProductConfigUpdateLogic) ProductConfigUpdate(in *dm.ProductConfig) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	old, err := relationDB.NewProductConfigRepo(l.ctx).FindOne(l.ctx, in.ProductID)
	if err != nil {
		return nil, err
	}
	po := utils.Copy[relationDB.DmProductConfig](in)
	old.DevInit = po.DevInit
	err = relationDB.NewProductConfigRepo(l.ctx).Update(l.ctx, old)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.ProductCache.SetData(l.ctx, in.ProductID, nil)
	if err != nil {
		l.Error(err)
	}
	err = l.svcCtx.FastEvent.Publish(l.ctx, topics.DmProductInfoUpdate, in.ProductID)
	if err != nil {
		l.Error(err)
	}
	return &dm.Empty{}, err
}
