package logic

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/share/devices"
	"github.com/zeromicro/go-zero/core/logx"
)

func UpdateDevice(ctx context.Context, svcCtx *svc.ServiceContext, devs []*devices.Core, affiliation devices.Affiliation) error {
	svcCtx.AbnormalRepo.UpdateDevice(ctx, devs, affiliation)
	svcCtx.SendRepo.UpdateDevice(ctx, devs, affiliation)
	for _, dev := range devs {
		s, err := svcCtx.DeviceSchemaRepo.GetData(ctx, *dev)
		if err != nil {
			logx.WithContext(ctx).Error(err.Error())
			continue
		}
		svcCtx.SchemaManaRepo.UpdateDevice(ctx, *dev, s, affiliation)
	}

	return nil
}
func UpdateDevices(ctx context.Context, svcCtx *svc.ServiceContext, devs []*devices.Info) error {
	svcCtx.AbnormalRepo.UpdateDevices(ctx, devs)
	svcCtx.SendRepo.UpdateDevices(ctx, devs)
	for _, dev := range devs {
		d := devices.Core{ProductID: dev.ProductID, DeviceName: dev.DeviceName}
		s, err := svcCtx.DeviceSchemaRepo.GetData(ctx, d)
		if err != nil {
			logx.WithContext(ctx).Error(err.Error())
			continue
		}
		svcCtx.SchemaManaRepo.UpdateDevice(ctx, d, s, utils.Copy2[devices.Affiliation](dev))
	}
	return nil
}

func UpdateDevGroupsTags(ctx context.Context, svcCtx *svc.ServiceContext, devs []devices.Core) error {
	var dgs []*devices.Info
	for _, g := range devs {
		gs, err := relationDB.NewGroupDeviceRepo(ctx).FindByFilter(ctx, relationDB.GroupDeviceFilter{WithGroup: true, ProductID: g.ProductID, DeviceName: g.DeviceName}, nil)
		if err != nil {
			logx.WithContext(ctx).Errorf("find group device error:%v", err)
			continue
		}
		var gp = map[string]def.IDsInfo{}
		for _, g := range gs {
			if g.Group == nil {
				continue
			}
			gg := gp[g.Group.Purpose]
			gg.IDs = append(gg.IDs, g.GroupID)
			gg.IDPaths = append(gg.IDPaths, g.Group.IDPath)
			gp[g.Group.Purpose] = gg
		}
		for i := 3; i > 0; i-- {
			dev := devices.Core{ProductID: g.ProductID, DeviceName: g.DeviceName}
			err := relationDB.NewDeviceInfoRepo(ctx).UpdateWithField(ctx, relationDB.DeviceFilter{Device: &dev}, map[string]any{
				"belong_group": utils.MarshalNoErr(gp),
			})
			if err != nil {
				logx.WithContext(ctx).Errorf("update group device error:%v", err)
				continue
			}
			svcCtx.DeviceCache.SetData(ctx, dev, nil)
			break
		}

		dgs = append(dgs, &devices.Info{
			ProductID:   g.ProductID,
			DeviceName:  g.DeviceName,
			BelongGroup: gp,
		})
	}
	return UpdateDevices(ctx, svcCtx, dgs)
}
