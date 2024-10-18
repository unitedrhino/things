package area

import (
	"context"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.SlotAreaSaveReq) error {
	//if req.ParentAreaID != def.RootNode {
	//	dmRep, err := l.svcCtx.DeviceM.DeviceInfoIndex(ctxs.WithRoot(l.ctx), &dm.DeviceInfoIndexReq{
	//		Page:    &dm.PageInfo{Page: 1, Size: 2}, //只需要知道是否有设备即可
	//		AreaIDs: []int64{req.ParentAreaID}})
	//	if err != nil {
	//		return err
	//	}
	//	if len(dmRep.List) != 0 {
	//		return errors.Parameter.AddMsg("父级区域已绑定了设备，不允许再添加子区域")
	//	}
	//}
	return nil
}
