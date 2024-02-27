package devicemanagelogic

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/domain/deviceMsg/msgGateway"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
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
func (l *DeviceGatewayMultiDeleteLogic) DeviceGatewayMultiDelete(in *dm.DeviceGatewayMultiDeleteReq) (*dm.Empty, error) {
	devicesDos := ToDeviceCoreDos(in.List)
	err := l.GdDB.MultiDelete(l.ctx, &devices.Core{
		ProductID:  in.GatewayProductID,
		DeviceName: in.GatewayDeviceName,
	}, devicesDos)
	if err != nil {
		return nil, err
	}
	req := &msgGateway.Msg{
		CommonMsg: deviceMsg.NewRespCommonMsg(l.ctx, deviceMsg.Change, "").AddStatus(errors.OK),
		Payload:   ToGatewayPayload(def.GatewayUnbind, devicesDos),
	}
	respBytes, _ := json.Marshal(req)
	msg := deviceMsg.PublishMsg{
		Handle:     devices.Gateway,
		Type:       msgGateway.TypeOperation,
		Payload:    respBytes,
		Timestamp:  time.Now().UnixMilli(),
		ProductID:  in.GatewayProductID,
		DeviceName: in.GatewayDeviceName,
	}
	er := l.svcCtx.PubDev.PublishToDev(l.ctx, &msg)
	if er != nil {
		l.Errorf("%s.PublishToDev failure err:%v", utils.FuncName(), er)
		return nil, er
	}
	return &dm.Empty{}, nil
}
