package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/dataType"
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
	"time"

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
		changeAreaIDPaths[string(di.AreaIDPath)] = struct{}{}
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
			changeAreaIDPaths[AreaIDPath] = struct{}{}
			projectIDSet[ai.ProjectID] = struct{}{}
		}
		pi, err = l.svcCtx.ProjectCache.GetData(l.ctx, in.ProjectID)
		if err != nil {
			return nil, err
		}
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
		err = relationDB.NewDeviceInfoRepo(tx).UpdateWithField(ctx, relationDB.DeviceFilter{Cores: devs}, param)
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
	for _, di := range devs {
		err = l.svcCtx.DeviceCache.SetData(l.ctx, *di, nil)
		if err != nil {
			l.Error(err)
		}
		err = l.svcCtx.FastEvent.Publish(l.ctx, topics.DmDeviceInfoUnbind, &di)
		if err != nil {
			l.Error(err)
		}
		BindChange(l.ctx, l.svcCtx, nil, *di, int64(ProjectID))

	}
	if len(changeAreaIDPaths) > 0 {
		ctxs.GoNewCtx(l.ctx, func(ctx2 context.Context) {
			logic.FillAreaDeviceCount(ctx2, l.svcCtx, utils.SetToSlice(changeAreaIDPaths)...)
		})
	}
	if len(projectIDSet) > 0 {
		ctxs.GoNewCtx(l.ctx, func(ctx2 context.Context) {
			logic.FillProjectDeviceCount(ctx2, l.svcCtx, utils.SetToSlice(projectIDSet)...)
		})
	}
	return &dm.Empty{}, nil
}
