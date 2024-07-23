package things

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/logic"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/i-Things/things/service/udsvr/pb/ud"
	"github.com/zeromicro/go-zero/core/logx"
)

func InfoToApi(ctx context.Context, svcCtx *svc.ServiceContext, v *dm.DeviceInfo, withProperties []string, withProfiles []string, withOwner bool) *types.DeviceInfo {
	var properties map[string]*types.DeviceInfoWithProperty
	var profiles map[string]string
	var owner *types.UserCore
	if err := ctxs.IsAdmin(ctx); err != nil {
		v.Secret = ""
		v.Cert = ""
	}
	position := &types.Point{
		Longitude: v.Position.Longitude, //经度
		Latitude:  v.Position.Latitude,  //维度
	}
	if withOwner && v.UserID > 0 {
		ui, err := svcCtx.UserM.UserInfoRead(ctxs.WithRoot(ctx), &sys.UserInfoReadReq{UserID: v.UserID})
		if err != nil {
			logx.WithContext(ctx).Error(err)
		} else {
			owner = utils.Copy[types.UserCore](ui)
		}
	}
	if withProperties != nil {
		func() {
			resp, err := svcCtx.DeviceMsg.PropertyLogLatestIndex(ctx, &dm.PropertyLogLatestIndexReq{
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
				DataIDs:    withProperties,
			})
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.PropertyLatestIndex err:%v", utils.FuncName(), err)
				return
			}
			if len(resp.List) == 0 {
				return
			}
			properties = make(map[string]*types.DeviceInfoWithProperty, len(resp.List))
			for _, v := range resp.List {
				properties[v.DataID] = &types.DeviceInfoWithProperty{
					Value:     v.Value,
					Timestamp: v.Timestamp,
				}
			}
		}()
	}
	if withProfiles != nil {
		ret, err := svcCtx.DeviceM.DeviceProfileIndex(ctx, &dm.DeviceProfileIndexReq{
			Device: &dm.DeviceCore{
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
			},
			Codes: withProfiles,
		})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.DeviceProfileIndex err:%v", utils.FuncName(), err)
		} else if len(ret.Profiles) > 0 {
			profiles = make(map[string]string, len(ret.Profiles))
			for _, v := range ret.Profiles {
				profiles[v.Code] = v.Params
			}
		}
	}
	if uc := ctxs.GetUserCtx(ctx); uc != nil && !uc.IsAdmin {
		v.Secret = "" // 设备秘钥
		v.Cert = ""   // 设备证书
	}
	return &types.DeviceInfo{
		ID:                 v.Id,
		TenantCode:         v.TenantCode,
		ProductID:          v.ProductID,          //产品id 只读
		DeviceName:         v.DeviceName,         //设备名称 读写
		DeviceAlias:        &v.DeviceAlias.Value, //设备别名 读写
		Secret:             v.Secret,             //设备秘钥 只读
		Cert:               v.Cert,               //设备证书 只读
		IsEnable:           v.IsEnable,
		Imei:               v.Imei,                        //IMEI号信息 只读
		Mac:                v.Mac,                         //MAC号信息 只读
		Version:            utils.ToNullString(v.Version), //固件版本 读写
		Rssi:               utils.ToEmptyInt64(v.Rssi),
		HardInfo:           v.HardInfo,               //模组硬件型号 只读
		SoftInfo:           v.SoftInfo,               //模组软件版本 只读
		Position:           position,                 //设别定位（百度坐标）
		Address:            &v.Address.Value,         //详细地址
		Tags:               logic.ToTagsType(v.Tags), //设备标签
		ProtocolConf:       logic.ToTagsType(v.Tags),
		SchemaAlias:        v.SchemaAlias, //设备物模型别名,如果是结构体类型则key为xxx.xxx
		IsOnline:           v.IsOnline,    //在线状态 1离线 2在线 只读
		FirstBind:          v.FirstBind,
		FirstLogin:         v.FirstLogin,  //激活时间 只读
		LastLogin:          v.LastLogin,   //最后上线时间 只读
		LogLevel:           v.LogLevel,    //日志级别 1)关闭 2)错误 3)告警 4)信息 5)调试  读写
		CreatedTime:        v.CreatedTime, //创建时间 只读
		MobileOperator:     v.MobileOperator,
		Phone:              utils.ToNullString(v.Phone),
		Iccid:              utils.ToNullString(v.Iccid),
		ProjectID:          v.ProjectID, //项目id 只读
		AreaID:             v.AreaID,    //项目区域id 只读
		WithProperties:     properties,
		Profiles:           profiles,
		Status:             v.Status,
		Manufacturer:       utils.Copy[types.ManufacturerInfo](v.Manufacturer),
		Owner:              owner,
		ProductName:        v.ProductName,
		DeviceType:         v.DeviceType,
		NetType:            v.NetType,
		ExpTime:            utils.ToEmptyInt64(v.ExpTime),
		RatedPower:         v.RatedPower,
		NeedConfirmVersion: v.NeedConfirmVersion,
		ProductImg:         v.ProductImg,
		CategoryID:         v.CategoryID,
		UserID:             v.UserID,
		Distributor:        utils.Copy[types.IDPath](v.Distributor),
	}
}

func ToDmDeviceCorePb(in *types.DeviceCore) *dm.DeviceCore {
	if in == nil {
		return nil
	}
	return &dm.DeviceCore{
		DeviceName: in.DeviceName,
		ProductID:  in.ProductID,
	}
}

func ToDmDeviceCoresPb(in []*types.DeviceCore) []*dm.DeviceCore {
	if in == nil {
		return nil
	}
	var ret []*dm.DeviceCore
	for _, v := range in {
		ret = append(ret, &dm.DeviceCore{
			DeviceName: v.DeviceName,
			ProductID:  v.ProductID,
		})
	}
	return ret
}

func UdToDeviceCoreTypes(in *ud.DeviceCore) *types.DeviceCore {
	if in == nil {
		return nil
	}
	return &types.DeviceCore{
		DeviceName: in.DeviceName,
		ProductID:  in.ProductID,
	}
}

func ToUdDeviceCorePb(in *types.DeviceCore) *ud.DeviceCore {
	if in == nil {
		return nil
	}
	return &ud.DeviceCore{
		DeviceName: in.DeviceName,
		ProductID:  in.ProductID,
	}
}

func ToUdDeviceCoresPb(in []*types.DeviceCore) []*ud.DeviceCore {
	if in == nil {
		return nil
	}
	var ret []*ud.DeviceCore
	for _, v := range in {
		ret = append(ret, &ud.DeviceCore{
			DeviceName: v.DeviceName,
			ProductID:  v.ProductID,
		})
	}
	return ret
}
