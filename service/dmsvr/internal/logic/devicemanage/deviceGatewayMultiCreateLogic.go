package devicemanagelogic

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/deviceAuth"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/domain/deviceMsg/msgGateway"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
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
func (l *DeviceGatewayMultiCreateLogic) DeviceGatewayMultiCreate(in *dm.DeviceGatewayMultiCreateReq) (*dm.Empty, error) {
	pi, err := l.svcCtx.ProductCache.GetData(l.ctx, in.Gateway.ProductID)
	if err != nil {
		return nil, err
	}
	_, err = FilterCanBindSubDevices(l.ctx, l.svcCtx, &devices.Core{
		ProductID:  in.Gateway.ProductID,
		DeviceName: in.Gateway.DeviceName,
	}, utils.ToSliceWithFunc(in.List, func(in *dm.DeviceGatewayBindDevice) *devices.Core {
		return &devices.Core{
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
		}
	}), CheckDeviceExist|CheckDeviceType|CheckDeviceStrict)
	if err != nil {
		return nil, err
	}
	if in.IsAuthSign { //秘钥检查
		for _, device := range in.List {
			di, err := l.DiDB.FindOneByFilter(l.ctx, relationDB.DeviceFilter{ProductID: device.ProductID, DeviceNames: []string{device.DeviceName}})
			if err != nil { //检查是否找到
				return nil, err
			}
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

	devicesDos := logic.BindToDeviceCoreDos(in.List)
	err = l.GdDB.MultiInsert(l.ctx, &devices.Core{
		ProductID:  in.Gateway.ProductID,
		DeviceName: in.Gateway.DeviceName,
	}, devicesDos)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	req := &msgGateway.Msg{
		CommonMsg: *deviceMsg.NewRespCommonMsg(l.ctx, deviceMsg.Change, "").AddStatus(errors.OK),
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

type CheckDevice int64

const (
	CheckDeviceExist CheckDevice = 1 << iota
	CheckDeviceType
	CheckDeviceStrict //严格模式
)

func FilterCanBindSubDevices(ctx context.Context, svcCtx *svc.ServiceContext, gateway *devices.Core, subDevices []*devices.Core, checkDevice CheckDevice) (ret []*devices.Core, err error) {
	{ //检查是否是网关类型
		pi, err := relationDB.NewProductInfoRepo(ctx).FindOneByFilter(ctx, relationDB.ProductFilter{ProductIDs: []string{gateway.ProductID}})
		if err != nil {
			if errors.Cmp(err, errors.NotFind) {
				return nil, errors.Parameter.AddDetail("not find GatewayProductID id:" + cast.ToString(gateway.ProductID))
			}
			return nil, errors.Database.AddDetail(err)
		}
		if pi.DeviceType != def.DeviceTypeGateway {
			return nil, errors.Parameter.AddMsg("子设备只能由网关设备进行绑定")
		}
	}

	if checkDevice&CheckDeviceType == CheckDeviceType { //检查设备是否都是子设备类型
		var (
			deviceProductList []string
			deviceProductMap  = map[string]struct{}{}
		)
		for _, v := range subDevices {
			deviceProductMap[v.ProductID] = struct{}{}
		}
		deviceProductList = utils.SetToSlice(deviceProductMap)
		products, err := relationDB.NewProductInfoRepo(ctx).FindByFilter(ctx, relationDB.ProductFilter{
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
	for _, subDevice := range subDevices { //检查是否有子设备绑定了其他网关
		if checkDevice&CheckDeviceExist == CheckDeviceExist { //检查设备是否都是子设备类型
			_, err := svcCtx.DeviceCache.GetData(ctx, *subDevice)
			if err != nil {
				if checkDevice&CheckDeviceStrict == CheckDeviceStrict {
					return nil, err
				}
				continue
			}
		}
		c, err := relationDB.NewGatewayDeviceRepo(ctx).FindOneByFilter(ctx, relationDB.GatewayDeviceFilter{
			SubDevice: &devices.Core{
				ProductID:  subDevice.ProductID,
				DeviceName: subDevice.DeviceName,
			}})
		if err == nil && !(c.GatewayProductID == gateway.ProductID && c.GatewayDeviceName == gateway.DeviceName) { //绑定了其他设备
			if checkDevice&CheckDeviceStrict == CheckDeviceStrict {
				return nil, err
			}
			continue
		} else {
			if err != nil && !errors.Cmp(err, errors.NotFind) {
				if checkDevice&CheckDeviceStrict == CheckDeviceStrict {
					return nil, err
				}
				continue
			}
		}
		//未绑定或就是该网关绑定的
		ret = append(ret, subDevice)
	}
	return ret, nil

}
