package devicemanagelogic

import (
	"context"
	"time"

	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/dataType"
	shareEvents "gitee.com/unitedrhino/share/events"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/topics"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceTransferLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceTransferLogic {
	return &DeviceTransferLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

const (
	DeviceTransferToUser    = 1
	DeviceTransferToProject = 2
)

func Transfer() {

}

func (l *DeviceTransferLogic) DeviceTransfer(in *dm.DeviceTransferReq) (*dm.Empty, error) {
	uc := ctxs.GetUserCtx(l.ctx)
	if in.SrcProjectID != 0 {
		uc.ProjectID = in.SrcProjectID
		if !uc.IsAdmin {
			if uc.ProjectAuth == nil || uc.ProjectAuth[in.ProjectID] == nil || uc.ProjectAuth[in.ProjectID].AuthType != def.AuthAdmin {
				return nil, errors.Permissions.AddMsg("只有项目管理员才能创建区域")
			}
		}
	}
	diDB := relationDB.NewDeviceInfoRepo(l.ctx)
	var dis []*relationDB.DmDeviceInfo
	var changeAreas = []*sys.AreaInfo{}
	var projectIDSet = map[int64]struct{}{}
	var changeTenantCode bool
	var oldTenantCode string
	if in.Device != nil && in.Device.ProductID != "" {
		di, err := diDB.FindOneByFilter(ctxs.WithDefaultRoot(l.ctx), relationDB.DeviceFilter{
			ProductID:   in.Device.ProductID,
			DeviceNames: []string{in.Device.DeviceName},
		})
		if err != nil {
			return nil, err
		}
		dis = append(dis, di)
	}
	if len(in.Devices) != 0 {
		di, err := diDB.FindByFilter(l.ctx, relationDB.DeviceFilter{
			Cores: utils.CopySlice[devices.Core](in.Devices),
		}, nil)
		if err != nil {
			return nil, err
		}
		dis = append(dis, di...)
	}
	if len(dis) == 0 {
		return &dm.Empty{}, nil
	}
	for _, di := range dis {
		if di.ProjectID <= def.NotClassified && uc.IsAdmin {
			continue
		}
		if oldTenantCode == "" {
			oldTenantCode = string(di.TenantCode)
		}
		if oldTenantCode != string(di.TenantCode) {
			changeTenantCode = true
		}
		if di.AreaID > def.NotClassified {
			ai, err := l.svcCtx.AreaCache.GetData(l.ctx, int64(di.AreaID))
			if err != nil {
				l.Error(err)
			} else {
				changeAreas = append(changeAreas, ai)
			}
		}
		projectIDSet[int64(di.ProjectID)] = struct{}{}

	}
	var (
		ProjectID  dataType.ProjectID
		pi         *sys.ProjectInfo
		err        error
		AreaID     dataType.AreaID = def.NotClassified
		AreaIDPath string          = def.NotClassifiedPath
		UserID     int64
	)

	switch in.TransferTo {
	case DeviceTransferToUser:
		dp, err := l.svcCtx.DataM.DataProjectIndex(ctxs.WithAllProject(l.ctx), &sys.DataProjectIndexReq{
			Page: &sys.PageInfo{
				Page: 1,
				Size: 1,
				Orders: []*sys.PageInfo_OrderBy{{
					Field: "createdTime", //第一个一定是默认的
					Sort:  stores.OrderAsc,
				}},
			},
			AuthType:   def.AuthAdmin,
			TargetID:   in.UserID,
			TargetType: "user",
		})
		if err != nil {
			return nil, err
		}
		if len(dp.List) == 0 {
			return nil, errors.NotFind.AddMsg("用户未找到")
		}
		ProjectID = dataType.ProjectID(dp.List[0].ProjectID)
		pi, err = l.svcCtx.ProjectCache.GetData(l.ctx, dp.List[0].ProjectID)
		if err != nil {
			return nil, err
		}
		UserID = in.UserID
	case DeviceTransferToProject:
		ProjectID = dataType.ProjectID(in.ProjectID)
		pi, err = l.svcCtx.ProjectCache.GetData(l.ctx, in.ProjectID)
		if err != nil {
			return nil, err
		}
		if in.AreaID > def.NotClassified {
			ai, err := l.svcCtx.AreaCache.GetData(l.ctx, in.AreaID)
			if err != nil {
				return nil, err
			}
			if ai.ProjectID != in.ProjectID {
				return nil, errors.Parameter.AddMsg("项目不对")
			}
			AreaID = dataType.AreaID(ai.AreaID)
			AreaIDPath = ai.AreaIDPath
			changeAreas = append(changeAreas, ai)
		}
		projectIDSet[in.ProjectID] = struct{}{}
		if ctxs.IsRoot(l.ctx) != nil && pi.TenantCode != uc.TenantCode {
			return nil, errors.Permissions.AddMsg("非超管不能转移到其他租户")
		}
		UserID = pi.AdminUserID
	default:
		return nil, errors.Parameter.AddMsgf("transferTo not supprt:%v", in.TransferTo)
	}
	if ctxs.IsRoot(l.ctx) != nil && (pi.TenantCode != oldTenantCode || changeTenantCode == true) {
		return nil, errors.Permissions.AddMsg("非超管不能转移到其他租户")
	}
	if in.IsCleanData == def.True {
		for _, di := range dis {
			err := DeleteDeviceTimeData(l.ctx, l.svcCtx, di.ProductID, di.DeviceName, DeleteModeThing)
			if err != nil {
				return nil, err
			}
		}
	}
	var devs = utils.CopySlice[devices.Core](dis)
	transferInfos := make([]*shareEvents.DeviceTransferInfo, 0, len(dis))
	for _, di := range dis {
		transferInfos = append(transferInfos, &shareEvents.DeviceTransferInfo{
			ProductID:     di.ProductID,
			DeviceName:    di.DeviceName,
			OldTenantCode: string(di.TenantCode),
			OldProjectID:  int64(di.ProjectID),
			OldAreaID:     int64(di.AreaID),
		})
	}

	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewUserDeviceShareRepo(tx).DeleteByFilter(ctxs.WithRoot(l.ctx), relationDB.UserDeviceShareFilter{
			Devices: devs,
		})
		if err != nil {
			return err
		}
		if in.IsCleanData == def.True {
			err = relationDB.NewDeviceProfileRepo(tx).DeleteByFilter(ctxs.WithRoot(l.ctx),
				relationDB.DeviceProfileFilter{Devices: devs})
			if err != nil {
				return err
			}
		}
		ctx := ctxs.WithAllProject(l.ctx)
		var param = map[string]any{
			"project_id":   ProjectID,
			"user_id":      UserID,
			"area_id":      AreaID,
			"area_id_path": AreaIDPath,
		}
		var tc string
		if pi.TenantCode != oldTenantCode || changeTenantCode == true {
			ctx = ctxs.WithRoot(l.ctx)
			param["tenant_code"] = pi.TenantCode
			tc = pi.TenantCode
		}
		if in.IsCleanData == def.True {
			param["last_bind"] = time.Now()
		}
		err = relationDB.NewDeviceInfoRepo(tx).UpdateWithField(ctxs.WithRoot(ctx), relationDB.DeviceFilter{Cores: devs}, param)
		if err != nil {
			return err
		}
		logic.UpdateDevice(l.ctx, l.svcCtx, devs, devices.Affiliation{
			TenantCode: tc,
			ProjectID:  int64(ProjectID),
			AreaID:     int64(AreaID),
			AreaIDPath: AreaIDPath,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	// 使用 root 权限查询，避免跨租户转移后因上下文租户不匹配导致查询失败
	postQueryCtx := ctxs.WithRoot(ctxs.WithAllProject(l.ctx))
	diDB = relationDB.NewDeviceInfoRepo(postQueryCtx)
	for i, di := range devs {
		err = l.svcCtx.DeviceCache.SetData(l.ctx, *di, nil)
		if err != nil {
			l.Error(err)
		}
		newDi, err := diDB.FindOneByFilter(postQueryCtx, relationDB.DeviceFilter{
			ProductID:   di.ProductID,
			DeviceNames: []string{di.DeviceName},
		})
		if err != nil {
			// 转让事务已提交成功，此处仅查询新设备信息用于事件发布，不影响转让结果
			l.Errorf("DeviceTransfer 转让后查询设备失败 productID=%s deviceName=%s err=%v", di.ProductID, di.DeviceName, err)
		} else {
			transferInfos[i].NewTenantCode = string(newDi.TenantCode)
			transferInfos[i].NewProjectID = int64(newDi.ProjectID)
			transferInfos[i].NewAreaID = int64(newDi.AreaID)
			transferInfos[i].NewAreaIDPath = string(newDi.AreaIDPath)
		}
		if in.IsCleanData == def.True {
			err = l.svcCtx.FastEvent.Publish(l.ctx, topics.DmDeviceInfoUnbind, &di)
			if err != nil {
				l.Error(err)
			}
		} else {
			err = l.svcCtx.FastEvent.Publish(l.ctx, topics.DmDeviceTransfer, transferInfos[i])
			if err != nil {
				l.Error(err)
			}
		}
		BindChange(l.ctx, l.svcCtx, nil, *di, int64(ProjectID))

	}
	if len(changeAreas) > 0 {
		ctxs.GoNewCtx(l.ctx, func(ctx2 context.Context) {
			logic.FillAreaDeviceCount(ctx2, l.svcCtx, changeAreas...)
		})
	}
	if len(projectIDSet) > 0 {
		ctxs.GoNewCtx(l.ctx, func(ctx2 context.Context) {
			logic.FillProjectDeviceCount(ctx2, l.svcCtx, utils.SetToSlice(projectIDSet)...)
		})
	}
	return &dm.Empty{}, nil
}
