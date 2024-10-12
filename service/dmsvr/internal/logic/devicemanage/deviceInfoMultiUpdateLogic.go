package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"time"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceInfoMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoMultiUpdateLogic {
	return &DeviceInfoMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量更新设备状态
func (l *DeviceInfoMultiUpdateLogic) DeviceInfoMultiUpdate(in *dm.DeviceInfoMultiUpdateReq) (*dm.Empty, error) {
	if len(in.Devices) == 0 {
		return &dm.Empty{}, nil
	}
	if in.AreaID == def.RootNode {
		return nil, errors.Parameter.AddMsgf("设备不能在root节点的区域下")
	}

	var columns []string
	var Distributor stores.IDPathWithUpdate
	var areaIDPath string
	var projectIDSet = map[int64]struct{}{}
	var changeAreaIDPaths = map[string]struct{}{}
	var devs = logic.ToDeviceCores(in.Devices)
	if in.AreaID != 0 {
		columns = append(columns, "area_id", "area_id_path")
		ai, err := l.svcCtx.AreaCache.GetData(l.ctx, in.AreaID)
		if err != nil {
			return nil, err
		}
		areaIDPath = ai.AreaIDPath
		changeAreaIDPaths[areaIDPath] = struct{}{}
		for _, dev := range devs {
			val, err := l.svcCtx.DeviceCache.GetData(l.ctx, *dev)
			if err != nil {
				continue
			}
			changeAreaIDPaths[val.AreaIDPath] = struct{}{}
			projectIDSet[val.ProjectID] = struct{}{}
		}
	}
	if in.Distributor != nil {
		columns = append(columns, "distributor_id", "distributor_id_path", "distributor_updated_time")
		Distributor = utils.Copy2[stores.IDPathWithUpdate](in.Distributor)
		Distributor.UpdatedTime = time.Now()
	}
	if in.RatedPower != 0 {
		columns = append(columns, "rated_power")
	}
	err := relationDB.NewDeviceInfoRepo(l.ctx).MultiUpdate(l.ctx, devs,
		&relationDB.DmDeviceInfo{RatedPower: in.RatedPower, AreaID: stores.AreaID(in.AreaID), AreaIDPath: areaIDPath, Distributor: utils.Copy2[stores.IDPathWithUpdate](in.Distributor)}, columns...)
	if err != nil {
		return nil, err
	}
	for _, v := range in.Devices {
		err := l.svcCtx.DeviceCache.SetData(l.ctx, devices.Core{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		}, nil)
		if err != nil {
			l.Error(err)
		}
	}
	if len(changeAreaIDPaths) > 0 {
		ctxs.GoNewCtx(l.ctx, func(ctx2 context.Context) {
			logic.FillAreaDeviceCount(ctx2, l.svcCtx, utils.SetToSlice(changeAreaIDPaths)...)
			logic.FillProjectDeviceCount(ctx2, l.svcCtx, utils.SetToSlice(projectIDSet)...)
		})
	}
	return &dm.Empty{}, err
}
