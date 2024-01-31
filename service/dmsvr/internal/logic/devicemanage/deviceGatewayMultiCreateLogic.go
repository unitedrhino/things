package devicemanagelogic

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/devices"
	"gitee.com/i-Things/core/shared/domain/deviceAuth"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg/msgGateway"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/spf13/cast"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceGatewayMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
	DiDB *relationDB.DeviceInfoRepo
	GdDB *relationDB.GatewayDeviceRepo
}

func NewDeviceGatewayMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceGatewayMultiCreateLogic {
	return &DeviceGatewayMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
		GdDB:   relationDB.NewGatewayDeviceRepo(ctx),
	}
}

// 创建分组设备
func (l *DeviceGatewayMultiCreateLogic) DeviceGatewayMultiCreate(in *dm.DeviceGatewayMultiCreateReq) (*dm.Response, error) {
	{ //检查是否是网关类型
		pi, err := l.PiDB.FindOneByFilter(l.ctx, relationDB.ProductFilter{ProductIDs: []string{in.GatewayProductID}})
		if err != nil {
			if errors.Cmp(err, errors.NotFind) {
				return nil, errors.Parameter.AddDetail("not find GatewayProductID id:" + cast.ToString(in.GatewayProductID))
			}
			return nil, errors.Database.AddDetail(err)
		}
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
	for _, device := range in.List {
		di, err := l.DiDB.FindOneByFilter(l.ctx, relationDB.DeviceFilter{ProductID: device.ProductID, DeviceNames: []string{device.DeviceName}})
		if err != nil { //检查是否找到
			return nil, err
		}
		if in.IsAuthSign { //秘钥检查
			if device.Sign == nil {
				return nil, errors.Parameter.AddMsg("没有填写签名信息")
			}
			pi, err := deviceAuth.NewPwdInfo(device.Sign.Signature, device.Sign.SignMethod)
			if err != nil {
				return nil, err
			}
			sign := fmt.Sprintf("%v;%v;%v;%v", device.ProductID, device.DeviceName, device.Sign.Random, device.Sign.Timestamp)
			if err := pi.CmpPwd(sign, di.Secret); err != nil {
				return nil, err
			}
		}
	}
	devicesDos := BindToDeviceCoreDos(in.List)
	err := l.GdDB.MultiInsert(l.ctx, &devices.Core{
		ProductID:  in.GatewayProductID,
		DeviceName: in.GatewayDeviceName,
	}, devicesDos)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	req := &msgGateway.Msg{
		CommonMsg: deviceMsg.NewRespCommonMsg(l.ctx, deviceMsg.Change, "").AddStatus(errors.OK),
		Payload:   ToGatewayPayload(def.GatewayBind, devicesDos),
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
	return &dm.Response{}, nil
}
