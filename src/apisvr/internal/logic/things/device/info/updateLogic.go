package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.DeviceInfoUpdateReq) error {
	if req.Position != nil {
		//经度范围是0-180°，纬度范围是0-90°
		if req.Position.Longitude < 0 || req.Position.Longitude > 180 {
			l.Errorf("%s.rpc.ManageDevice req=%v err= Longitude value is invalid", utils.FuncName(), req)
			return errors.Parameter.AddDetail("Longitude value is invalid")
		}
		if req.Position.Latitude < 0 || req.Position.Latitude > 90 {
			l.Errorf("%s.rpc.ManageDevice req=%v err= the Latitude value is invalid", utils.FuncName(), req)
			return errors.Parameter.AddDetail("Latitude value is invalid")
		}
	}
	dmReq := &dm.DeviceInfo{
		ProductID:  req.ProductID,  //产品id 只读
		DeviceName: req.DeviceName, //设备名称 读写
		LogLevel:   req.LogLevel,   // 日志级别:1)关闭 2)错误 3)告警 4)信息 5)调试  读写
		Tags:       logic.ToTagsMap(req.Tags),
		Address:    utils.ToRpcNullString(req.Address),
		Position:   logic.ToDmPointRpc(req.Position),
	}
	_, err := l.svcCtx.DeviceM.DeviceInfoUpdate(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ManageDevice req=%v err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
