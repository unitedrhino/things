package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceStatus"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceAuth"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceGatewayMultiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	GdDB *relationDB.GatewayDeviceRepo
}

func NewDeviceGatewayMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceGatewayMultiDeleteLogic {
	return &DeviceGatewayMultiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		GdDB:   relationDB.NewGatewayDeviceRepo(ctx),
	}
}

// 删除分组设备
func (l *DeviceGatewayMultiDeleteLogic) DeviceGatewayMultiDelete(in *dm.DeviceGatewayMultiSaveReq) (*dm.Empty, error) {
	pi, err := l.svcCtx.ProductCache.GetData(l.ctx, in.Gateway.ProductID)
	if err != nil {
		return nil, err
	}
	devicesDos := logic.ToDeviceCoreDos(in.List)
	list, err := l.GdDB.FindByFilter(l.ctx, relationDB.GatewayDeviceFilter{
		Gateway: &devices.Core{
			ProductID:  in.Gateway.ProductID,
			DeviceName: in.Gateway.DeviceName,
		},
		SubDevices: devicesDos,
	}, nil)
	if err != nil {
		return nil, err
	}
	if len(list) != len(devicesDos) {
		return &dm.Empty{}, errors.Permissions.AddMsg("有子设备未挂载到该网关下")
	}
	_, err = NewDeviceInfoMultiUpdateLogic(ctxs.WithProjectID(l.ctx, def.NotClassified), l.svcCtx).DeviceInfoMultiUpdate(&dm.DeviceInfoMultiUpdateReq{
		Devices: in.List,
		AreaID:  def.NotClassified,
	})
	if err != nil {
		return nil, err
	}
	gateway := devices.Core{
		ProductID:  in.Gateway.ProductID,
		DeviceName: in.Gateway.DeviceName,
	}
	err = l.GdDB.MultiDelete(l.ctx, &gateway, devicesDos)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	for _, v := range devicesDos {
		//更新在线状态
		err := HandleOnlineFix(l.ctx, l.svcCtx, &deviceStatus.ConnectMsg{
			ClientID:  deviceAuth.GenClientID(v.ProductID, v.DeviceName),
			Timestamp: now,
			Action:    devices.ActionDisconnected,
			Reason:    "gateway unbind",
		})
		if err != nil {
			l.Error(err)
		}
	}
	if in.IsNotNotify {
		return &dm.Empty{}, nil
	}
	TopoChange(l.ctx, l.svcCtx, def.GatewayUnbind, pi, gateway, devicesDos)

	return &dm.Empty{}, nil
}
