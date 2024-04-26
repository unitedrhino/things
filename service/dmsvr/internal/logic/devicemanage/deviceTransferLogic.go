package devicemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
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

func (l *DeviceTransferLogic) DeviceTransfer(in *dm.DeviceTransferReq) (*dm.Empty, error) {
	diDB := relationDB.NewDeviceInfoRepo(l.ctx)
	di, err := diDB.FindOneByFilter(l.ctx, relationDB.DeviceFilter{
		ProductID:   in.Device.ProductID,
		DeviceNames: []string{in.Device.DeviceName},
	})
	if err != nil {
		return nil, err
	}
	pi, err := l.svcCtx.ProjectM.ProjectInfoRead(l.ctx, &sys.ProjectWithID{ProjectID: int64(di.ProjectID)})
	if err != nil {
		return nil, err
	}
	if pi.AdminUserID != pi.AdminUserID {
		return nil, errors.Permissions
	}
	switch in.TransferTo {
	case DeviceTransferToUser:
		dp, err := l.svcCtx.DataM.DataProjectIndex(l.ctx, &sys.DataProjectIndexReq{
			Page: &sys.PageInfo{
				Page: 1,
				Size: 1,
				Orders: []*sys.PageInfo_OrderBy{{
					Filed: "createdTime", //第一个一定是默认的
					Sort:  def.OrderAsc,
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
		di.ProjectID = stores.ProjectID(dp.List[0].ProjectID)
		di.AreaID = def.NotClassified
	case DeviceTransferToProject:
		di.ProjectID = stores.ProjectID(pi.ProjectID)
		di.AreaID = def.NotClassified
	}
	if in.IsCleanData == def.True {
		err = DeleteDeviceTimeData(l.ctx, l.svcCtx, in.Device.ProductID, in.Device.DeviceName)
		if err != nil {
			return nil, err
		}
	}
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewUserDeviceShareRepo(tx).DeleteByFilter(l.ctx, relationDB.UserDeviceShareFilter{
			ProductID:  in.Device.ProductID,
			DeviceName: in.Device.DeviceName,
		})
		if err != nil {
			return err
		}
		err = relationDB.NewDeviceInfoRepo(tx).Update(ctxs.WithAllProject(l.ctx), di)
		if err != nil {
			return err
		}
		return nil
	})

	return &dm.Empty{}, err
}
