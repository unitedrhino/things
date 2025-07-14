package things

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoWith struct {
	Properties []string
	Profiles   []string
	WithGroups []string
	Owner      bool
	Area       bool
	IsOnlyCore bool
}

func InfoToApi(ctx context.Context, svcCtx *svc.ServiceContext, v *dm.DeviceInfo, w DeviceInfoWith) *types.DeviceInfo {
	if v == nil {
		return nil
	}
	var properties map[string]*types.DeviceInfoWithProperty
	var profiles map[string]string
	var owner *types.UserCore
	var area *types.AreaInfo
	if err := ctxs.IsAdmin(ctx); err != nil {
		v.Secret = ""
		v.Cert = ""
	}
	position := &types.Point{
		Longitude: v.Position.Longitude, //经度
		Latitude:  v.Position.Latitude,  //维度
	}
	if w.Owner && v.UserID > 0 {
		ui, err := svcCtx.UserC.GetData(ctx, v.UserID)
		if err != nil {
			logx.WithContext(ctx).Error(v.UserID, err)
		} else {
			owner = utils.Copy[types.UserCore](ui)
		}
	}
	if w.Area {
		a, err := svcCtx.AreaC.GetData(ctx, v.AreaID)
		if err != nil {
			logx.WithContext(ctx).Error(err)
		} else {
			area = utils.Copy[types.AreaInfo](a)
		}
	}
	if w.Properties != nil {
		func() {
			resp, err := svcCtx.DeviceMsg.PropertyLogLatestIndex(ctx, &dm.PropertyLogLatestIndexReq{
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
				DataIDs:    w.Properties,
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
	if w.Profiles != nil {
		ret, err := svcCtx.DeviceM.DeviceProfileIndex(ctx, &dm.DeviceProfileIndexReq{
			Device: &dm.DeviceCore{
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
			},
			Codes: w.Profiles,
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
		v.FirstBind = 0
		v.FirstLogin = 0
	}
	var groups []*types.GroupCore
	if w.WithGroups != nil {
		gs, err := svcCtx.DeviceG.GroupInfoIndex(ctx, &dm.GroupInfoIndexReq{
			Purposes:  w.WithGroups,
			HasDevice: &dm.DeviceCore{ProductID: v.ProductID, DeviceName: v.DeviceName},
		})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.GroupInfoIndex err:%v", utils.FuncName(), err)
		} else {
			groups = utils.CopySlice[types.GroupCore](gs.List)
		}
	}
	if w.IsOnlyCore {
		return &types.DeviceInfo{
			ID:             v.Id,
			ProductID:      v.ProductID,          //产品id 只读
			DeviceName:     v.DeviceName,         //设备名称 读写
			DeviceAlias:    &v.DeviceAlias.Value, //设备别名 读写
			Rssi:           utils.ToEmptyInt64(v.Rssi),
			IsOnline:       v.IsOnline,  //在线状态 1离线 2在线 只读
			ProjectID:      v.ProjectID, //项目id 只读
			AreaID:         v.AreaID,    //项目区域id 只读
			WithProperties: properties,
			Status:         v.Status,
			ProductName:    v.ProductName,
			NetType:        v.NetType,
			ProductImg:     v.ProductImg,
			Distributor:    utils.Copy[types.IDPath](v.Distributor),
			Gateway:        InfoToApi(ctx, svcCtx, v.Gateway, DeviceInfoWith{IsOnlyCore: true}),
			Area:           area,
			Groups:         groups,
		}
	}
	return &types.DeviceInfo{
		ID:                 v.Id,
		TenantCode:         v.TenantCode,
		ProductID:          v.ProductID,          //产品id 只读
		DeviceName:         v.DeviceName,         //设备名称 读写
		DeviceAlias:        &v.DeviceAlias.Value, //设备别名 读写
		Secret:             v.Secret,             //设备秘钥 只读
		Cert:               v.Cert,               //设备证书 只读
		IsUpdateDeviceImg:  v.IsUpdateDeviceImg,
		IsUpdateFile:       v.IsUpdateFile,
		DeviceImg:          v.DeviceImg,
		File:               v.File,
		LastLocalIp:        v.LastLocalIp,
		LastOffline:        v.LastOffline,
		IsEnable:           v.IsEnable,
		Imei:               v.Imei,                        //IMEI号信息 只读
		Mac:                v.Mac,                         //MAC号信息 只读
		Version:            utils.ToNullString(v.Version), //固件版本 读写
		Rssi:               utils.ToEmptyInt64(v.Rssi),
		HardInfo:           v.HardInfo,                    //模组硬件型号 只读
		SoftInfo:           v.SoftInfo,                    //模组软件版本 只读
		Position:           position,                      //设别定位（百度坐标）
		Address:            utils.ToNullString(v.Address), //详细地址
		Adcode:             utils.ToNullString(v.Adcode),
		Tags:               logic.ToTagsType(v.Tags), //设备标签
		ProtocolConf:       logic.ToTagsType(v.ProtocolConf),
		SubProtocolConf:    logic.ToTagsType(v.SubProtocolConf),
		SchemaAlias:        v.SchemaAlias, //设备物模型别名,如果是结构体类型则key为xxx.xxx
		IsOnline:           v.IsOnline,    //在线状态 1离线 2在线 只读
		FirstBind:          v.FirstBind,
		LastBind:           v.LastBind,
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
		Owner:              owner,
		ProductName:        v.ProductName,
		DeviceType:         v.DeviceType,
		Sort:               v.Sort,
		NetType:            v.NetType,
		ExpTime:            utils.ToNullInt64(v.ExpTime),
		RatedPower:         v.RatedPower,
		NeedConfirmVersion: v.NeedConfirmVersion,
		ProductImg:         v.ProductImg,
		CategoryID:         v.CategoryID,
		UserID:             v.UserID,
		Desc:               utils.ToNullString(v.Desc),
		Distributor:        utils.Copy[types.IDPath](v.Distributor),
		Gateway:            InfoToApi(ctx, svcCtx, v.Gateway, DeviceInfoWith{IsOnlyCore: true}),
		Area:               area,
		LastIp:             v.LastIp,
		Groups:             groups,
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
