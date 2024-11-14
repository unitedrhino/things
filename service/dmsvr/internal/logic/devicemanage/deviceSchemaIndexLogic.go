package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceSchemaIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceSchemaIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceSchemaIndexLogic {
	return &DeviceSchemaIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备物模型列表
func (l *DeviceSchemaIndexLogic) DeviceSchemaIndex(in *dm.DeviceSchemaIndexReq) (*dm.DeviceSchemaIndexResp, error) {
	l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(in))
	filter := utils.Copy2[relationDB.ProductSchemaFilter](in)
	schemas, err := relationDB.NewProductSchemaRepo(l.ctx).FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page).WithDefaultOrder(stores.OrderBy{
		Field: "order",
		Sort:  stores.OrderAsc,
	}))
	if err != nil {
		return nil, err
	}
	total, err := relationDB.NewProductSchemaRepo(l.ctx).CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	list := utils.CopySlice[dm.DeviceSchema](schemas)
	dfilter := utils.Copy2[relationDB.DeviceSchemaFilter](in)
	total2, err := relationDB.NewDeviceSchemaRepo(l.ctx).CountByFilter(l.ctx, dfilter)
	if err != nil {
		return nil, err
	}
	if in.Page != nil && in.Page.Size <= int64(len(list)) { //如果分页获取的数据还没获取完需要从设备物模型中获取
		return &dm.DeviceSchemaIndexResp{List: list, Total: total + total2}, nil
	}

	//dfilter := utils.Copy2[relationDB.DeviceSchemaFilter](in)
	//schemas, err := relationDB.NewDeviceSchemaRepo(l.ctx).FindByFilter(l.ctx, dfilter, logic.ToPageInfo(in.Page).WithDefaultOrder(stores.OrderBy{
	//	Field: "order",
	//	Sort:  stores.OrderAsc,
	//}))
	//if err != nil {
	//	return nil, err
	//}
	//total, err := relationDB.NewDeviceSchemaRepo(l.ctx).CountByFilter(l.ctx, dfilter)
	//if err != nil {
	//	return nil, err
	//}
	return &dm.DeviceSchemaIndexResp{List: utils.CopySlice[dm.DeviceSchema](schemas), Total: total}, nil
}
