package devicemanagelogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgGateway"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

// BindChange 同一个项目只通知一次
func BindChange(ctx context.Context, svcCtx *svc.ServiceContext, pi *dm.ProductInfo, dev devices.Core, projectID int64) error {
	v := svcCtx.BindChange.Get(ctx, dev)
	if v == projectID {
		return nil
	}
	req := &msgThing.Req{
		CommonMsg: *deviceMsg.NewRespCommonMsg(ctx, deviceMsg.BindChange, devices.GenMsgToken(ctx, svcCtx.NodeID)).AddStatus(errors.OK, false),
		Params:    map[string]any{"projectID": cast.ToString(projectID)},
	}
	respBytes, _ := json.Marshal(req)
	if pi == nil {
		var err error
		pi, err = svcCtx.ProductCache.GetData(ctx, dev.ProductID)
		logx.WithContext(ctx).Error(dev.ProductID, err)
		if err != nil {
			return err
		}
	}
	msg := deviceMsg.PublishMsg{
		Handle:       devices.Thing,
		Type:         msgThing.TypeService,
		Payload:      respBytes,
		Timestamp:    time.Now().UnixMilli(),
		ProductID:    dev.ProductID,
		DeviceName:   dev.DeviceName,
		ProtocolCode: pi.ProtocolCode,
	}
	er := svcCtx.PubDev.PublishToDev(ctx, &msg)
	if er != nil {
		logx.WithContext(ctx).Errorf("%s.PublishToDev failure err:%v", utils.FuncName(), er)
		return er
	}
	svcCtx.BindChange.Set(ctx, dev, projectID)
	return nil
}

func TopoChange(ctx context.Context, svcCtx *svc.ServiceContext, status def.GatewayStatus, pi *dm.ProductInfo, gateway devices.Core, devs []*devices.Core) error {
	req := &msgGateway.Msg{
		CommonMsg: *deviceMsg.NewRespCommonMsg(ctx, deviceMsg.Change, devices.GenMsgToken(ctx, svcCtx.NodeID)).AddStatus(errors.OK, false),
		Payload:   logic.ToGatewayPayload(status, devs),
	}
	respBytes, _ := json.Marshal(req)
	msg := deviceMsg.PublishMsg{
		Handle:       devices.Gateway,
		Type:         msgGateway.TypeTopo,
		Payload:      respBytes,
		Timestamp:    time.Now().UnixMilli(),
		ProductID:    gateway.ProductID,
		DeviceName:   gateway.DeviceName,
		ProtocolCode: pi.ProtocolCode,
	}
	er := svcCtx.PubDev.PublishToDev(ctx, &msg)
	if er != nil {
		logx.WithContext(ctx).Errorf("%s.PublishToDev failure err:%v", utils.FuncName(), er)
		return er
	}
	return nil
}
