package devicemanagelogic

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceAuth"
	"github.com/spf13/cast"
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
	gd, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
		ProductID:  in.Gateway.ProductID,
		DeviceName: in.Gateway.DeviceName,
	})
	if err != nil {
		return nil, err
	}
	devs, err := FilterCanBindSubDevices(l.ctx, l.svcCtx, &devices.Core{
		ProductID:  in.Gateway.ProductID,
		DeviceName: in.Gateway.DeviceName,
	}, utils.ToSliceWithFunc(in.List, func(in *dm.DeviceGatewayBindDevice) *devices.Core {
		return &devices.Core{
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
		}
	}), CheckDeviceExist|CheckDeviceType)
	if err != nil {
		return nil, err
	}
	if len(devs) == 0 {
		return &dm.Empty{}, nil
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
	gateway := devices.Core{
		ProductID:  in.Gateway.ProductID,
		DeviceName: in.Gateway.DeviceName,
	}
	err = l.GdDB.MultiInsert(l.ctx, &gateway, devs)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	_, err = NewDeviceTransferLogic(ctxs.WithRoot(l.ctx), l.svcCtx).DeviceTransfer(&dm.DeviceTransferReq{
		TransferTo: DeviceTransferToProject,
		Devices:    utils.CopySlice[dm.DeviceCore](devs),
		AreaID:     gd.AreaID,
		ProjectID:  gd.ProjectID,
	})
	if err != nil {
		l.Error(err)
	}
	if in.IsNotNotify {
		return &dm.Empty{}, nil
	}
	TopoChange(l.ctx, l.svcCtx, def.GatewayBind, pi, gateway, devs)
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
		pi, err := svcCtx.ProductCache.GetData(ctx, gateway.ProductID)
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
	if len(subDevices) == 0 {
		return []*devices.Core{}, nil
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
		if checkDevice&CheckDeviceExist == CheckDeviceExist { //检查设备是否都存在
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
		if err == nil { //绑定了其他设备
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
		if c != nil && c.GatewayProductID == gateway.ProductID && c.GatewayDeviceName == gateway.DeviceName { //如果已经绑定了就忽略
			continue
		}
		//未绑定或就是该网关绑定的
		ret = append(ret, subDevice)
	}
	return ret, nil

}
