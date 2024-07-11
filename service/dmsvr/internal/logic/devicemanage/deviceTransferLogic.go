package devicemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"gorm.io/gorm"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

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
	diDB := relationDB.NewDeviceInfoRepo(l.ctx)
	var dis []*relationDB.DmDeviceInfo
	if in.Device != nil {
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
	}
	var (
		ProjectID  stores.ProjectID
		AreaID     stores.AreaID = def.NotClassified
		AreaIDPath string        = def.NotClassifiedPath
	)

	switch in.TransferTo {
	case DeviceTransferToUser:
		dp, err := l.svcCtx.DataM.DataProjectIndex(l.ctx, &sys.DataProjectIndexReq{
			Page: &sys.PageInfo{
				Page: 1,
				Size: 1,
				Orders: []*sys.PageInfo_OrderBy{{
					Field: "createdTime", //第一个一定是默认的
					Sort:  stores.OrderAsc,
				}},
			},
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
	case DeviceTransferToProject:
		ProjectID = stores.ProjectID(in.ProjectID)
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
	}

	return &dm.Empty{}, nil
}
