package devicegrouplogic

import (
	"context"
	"gitee.com/unitedrhino/share/oss/common"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
)

func ToGroupInfoPb(ctx context.Context, svcCtx *svc.ServiceContext, ro *relationDB.DmGroupInfo) *dm.GroupInfo {
	if ro == nil {
		return nil
	}
	productName := ""
	if ro.ProductInfo != nil {
		productName = ro.ProductInfo.ProductName
	}
	for k, v := range ro.Files {
		if v == "" {
			continue
		}
		var err error
		ro.Files[k], err = svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, v, 60*60, common.OptionKv{})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.SignedGetUrl err:%v", utils.FuncName(), err)
		}
	}

	return &dm.GroupInfo{
		AreaID:      int64(ro.AreaID),
		Id:          ro.ID,
		IdPath:      ro.IDPath,
		ParentID:    ro.ParentID,
		ProjectID:   int64(ro.ProjectID),
		ProductName: productName,
		Name:        ro.Name,
		Files:       ro.Files,
		ProductID:   ro.ProductID,
		DeviceCount: ro.DeviceCount,
		Desc:        ro.Desc,
		CreatedTime: ro.CreatedTime.Unix(),
		Purpose:     ro.Purpose,
		Tags:        ro.Tags,
		IsLeaf:      ro.IsLeaf,
	}
}

func fillGroupDevices(groups []*dm.GroupInfo, groupDevices []*relationDB.DmGroupDevice) {
	if len(groups) == 0 || len(groupDevices) == 0 {
		return
	}
	groupDeviceMap := make(map[int64][]*dm.DeviceCore)
	for _, gd := range groupDevices {
		if gd == nil || gd.Device == nil {
			continue
		}
		groupDeviceMap[gd.GroupID] = append(groupDeviceMap[gd.GroupID], &dm.DeviceCore{
			ProductID:  gd.Device.ProductID,
			DeviceName: gd.Device.DeviceName,
		})
	}
	fillGroupDevicesByMap(groups, groupDeviceMap)
}

func fillGroupDevicesByMap(groups []*dm.GroupInfo, groupDeviceMap map[int64][]*dm.DeviceCore) {
	for _, group := range groups {
		if group == nil {
			continue
		}
		if devices := groupDeviceMap[group.Id]; len(devices) > 0 {
			group.Devices = devices
		}
		fillGroupDevicesByMap(group.Children, groupDeviceMap)
	}
}
