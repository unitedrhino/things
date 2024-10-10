package devicemanagelogic

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/deviceAuth"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/domain/deviceMsg/msgGateway"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/dmsvr/internal/domain/deviceStatus"
	"gitee.com/i-Things/things/service/dmsvr/internal/logic"
	"gitee.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"
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
	_, err = NewDeviceInfoMultiUpdateLogic(ctxs.WithRoot(l.ctx), l.svcCtx).DeviceInfoMultiUpdate(&dm.DeviceInfoMultiUpdateReq{
		Devices: in.List,
		AreaID:  def.NotClassified,
	})
	if err != nil {
		return nil, err
	}

	err = l.GdDB.MultiDelete(l.ctx, &devices.Core{
		ProductID:  in.Gateway.ProductID,
		DeviceName: in.Gateway.DeviceName,
	}, devicesDos)
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
	req := &msgGateway.Msg{
		CommonMsg: *deviceMsg.NewRespCommonMsg(l.ctx, deviceMsg.Change, devices.GenMsgToken(l.ctx, l.svcCtx.NodeID)).AddStatus(errors.OK),
		Payload:   logic.ToGatewayPayload(def.GatewayUnbind, devicesDos),
	}
	respBytes, _ := json.Marshal(req)
	msg := deviceMsg.PublishMsg{
		Handle:       devices.Gateway,
		Type:         msgGateway.TypeTopo,
		Payload:      respBytes,
		Timestamp:    now.UnixMilli(),
		ProductID:    in.Gateway.ProductID,
		DeviceName:   in.Gateway.DeviceName,
		ProtocolCode: pi.ProtocolCode,
	}
	er := l.svcCtx.PubDev.PublishToDev(l.ctx, &msg)
	if er != nil {
		l.Errorf("%s.PublishToDev failure err:%v", utils.FuncName(), er)
		return nil, er
	}

	return &dm.Empty{}, nil
}
