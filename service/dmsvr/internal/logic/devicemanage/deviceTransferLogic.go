package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
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
	var changeAreaIDPaths = map[string]struct{}{}
	var projectIDSet = map[int64]struct{}{}
	if in.Device != nil && in.Device.ProductID != "" {
		di, err := diDB.FindOneByFilter(l.ctx, relationDB.DeviceFilter{
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
		pi, err := l.svcCtx.ProjectCache.GetData(l.ctx, int64(di.ProjectID))
		if err != nil {
			return nil, err
		}
		if pi.AdminUserID != pi.AdminUserID {
			return nil, errors.Permissions
		}
		changeAreaIDPaths[di.AreaIDPath] = struct{}{}
		projectIDSet[int64(di.ProjectID)] = struct{}{}

	}
	var (
		ProjectID  stores.ProjectID
		AreaID     stores.AreaID = def.NotClassified
		AreaIDPath string        = def.NotClassifiedPath
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
		ProjectID = stores.ProjectID(dp.List[0].ProjectID)
		UserID = in.UserID
	case DeviceTransferToProject:
		ProjectID = stores.ProjectID(in.ProjectID)
		if in.AreaID != 0 {
			ai, err := l.svcCtx.AreaCache.GetData(l.ctx, in.AreaID)
			if err != nil {
				return nil, err
			}
			if ai.ProjectID != in.ProjectID {
				return nil, errors.Parameter.AddMsg("项目不对")
			}
			AreaID = stores.AreaID(ai.AreaID)
			AreaIDPath = ai.AreaIDPath
			changeAreaIDPaths[AreaIDPath] = struct{}{}
			projectIDSet[ai.ProjectID] = struct{}{}

		}
		pi, err := l.svcCtx.ProjectCache.GetData(l.ctx, in.ProjectID)
		if err != nil {
			return nil, err
		}
		UserID = pi.AdminUserID
	default:
		return nil, errors.Parameter.AddMsgf("transferTo not supprt:%v", in.TransferTo)
	}
	if in.IsCleanData == def.True {
		for _, di := range dis {
			err := DeleteDeviceTimeData(l.ctx, l.svcCtx, di.ProductID, di.DeviceName)
			if err != nil {
				return nil, err
			}
		}
	}
	var devs = utils.CopySlice[devices.Core](dis)
	err := stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewUserDeviceShareRepo(tx).DeleteByFilter(l.ctx, relationDB.UserDeviceShareFilter{
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
		err = relationDB.NewDeviceInfoRepo(tx).UpdateWithField(ctxs.WithAllProject(l.ctx), relationDB.DeviceFilter{Cores: devs}, map[string]any{
			"project_id":   ProjectID,
			"user_id":      UserID,
			"area_id":      AreaID,
			"area_id_path": AreaIDPath,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	for _, di := range devs {
		err = l.svcCtx.DeviceCache.SetData(l.ctx, *di, nil)
		if err != nil {
			l.Error(err)
		}
		err = l.svcCtx.FastEvent.Publish(l.ctx, eventBus.DmDeviceInfoUnbind, &di)
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
	return &dm.Empty{}, nil
}
