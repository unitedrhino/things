package devicemanagelogic

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/domain/deviceMsg/msgGateway"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"time"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceGatewayMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
	DiDB *relationDB.DeviceInfoRepo
	GdDB *relationDB.GatewayDeviceRepo
}

func NewDeviceGatewayMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceGatewayMultiUpdateLogic {
	return &DeviceGatewayMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
		GdDB:   relationDB.NewGatewayDeviceRepo(ctx),
	}
}

// 绑定网关下子设备设备
func (l *DeviceGatewayMultiUpdateLogic) DeviceGatewayMultiUpdate(in *dm.DeviceGatewayMultiSaveReq) (*dm.Empty, error) {
	pi, err := l.svcCtx.ProductCache.GetData(l.ctx, in.Gateway.ProductID)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddDetail("not find GatewayProductID id:" + cast.ToString(in.Gateway.ProductID))
		}
		return nil, errors.Database.AddDetail(err)
	}
	{ //检查是否是网关类型
		if pi.DeviceType != def.DeviceTypeGateway {
			return nil, errors.Parameter.AddMsg("子设备只能由网关设备进行绑定")
		}
	}
	{ //检查设备是否都是子设备类型
		var (
			deviceProductList []string
			deviceProductMap  = map[string]struct{}{}
		)
		for _, v := range in.List {
			deviceProductMap[v.ProductID] = struct{}{}
		}
		deviceProductList = utils.SetToSlice(deviceProductMap)
		products, err := l.PiDB.FindByFilter(l.ctx, relationDB.ProductFilter{
			ProductIDs: deviceProductList,
		}, nil)
		if err != nil {
			return nil, errors.Database.AddDetail(err)
		}
		for _, v := range products {
			if v.DeviceType != def.DeviceTypeSubset {
				return nil, errors.Parameter.AddMsg("网关只能绑定子设备类型")
			}
		}
	}
	devicesDos := logic.ToDeviceCoreDos(in.List)
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		gd := relationDB.NewGatewayDeviceRepo(tx)
		err := gd.MultiDelete(l.ctx, &devices.Core{
			ProductID:  in.Gateway.ProductID,
			DeviceName: in.Gateway.DeviceName,
		}, nil)
		if err != nil {
			return err
		}
		err = gd.MultiInsert(l.ctx, &devices.Core{
			ProductID:  in.Gateway.ProductID,
			DeviceName: in.Gateway.DeviceName,
		}, devicesDos)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	req := &msgGateway.Msg{
		CommonMsg: *deviceMsg.NewRespCommonMsg(l.ctx, deviceMsg.Change, devices.GenMsgToken(l.ctx, l.svcCtx.NodeID)).AddStatus(errors.OK),
		Payload:   logic.ToGatewayPayload(def.GatewayBind, devicesDos),
	}
	respBytes, _ := json.Marshal(req)
	msg := deviceMsg.PublishMsg{
		Handle:       devices.Gateway,
		Type:         msgGateway.TypeTopo,
		Payload:      respBytes,
		Timestamp:    time.Now().UnixMilli(),
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
