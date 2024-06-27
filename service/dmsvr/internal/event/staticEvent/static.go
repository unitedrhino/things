package staticEvent

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type StaticHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewStaticHandle(ctx context.Context, svcCtx *svc.ServiceContext) *StaticHandle {
	return &StaticHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (l *StaticHandle) Handle() error { //产品品类设备数量统计
	pcDB := relationDB.NewProductCategoryRepo(l.ctx)
	pcs, err := pcDB.FindByFilter(l.ctx, relationDB.ProductCategoryFilter{}, nil)
	if err != nil {
		return err
	}
	for _, pc := range pcs {
		ids := utils.GetIDPath(pc.IDPath)
		total, err := relationDB.NewDeviceInfoRepo(l.ctx).CountByFilter(l.ctx, relationDB.DeviceFilter{ProductCategoryIDs: ids})
		if err != nil {
			l.Error(err)
			continue
		}
		pc.DeviceCount = total
		err = pcDB.Update(l.ctx, pc)
		if err != nil {
			l.Error(err)
			continue
		}
	}
	return nil
}
