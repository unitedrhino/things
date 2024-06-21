package devicemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"gorm.io/gorm"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoUnbindLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceInfoUnbindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoUnbindLogic {
	return &DeviceInfoUnbindLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceInfoUnbindLogic) DeviceInfoUnbind(in *dm.DeviceCore) (*dm.Empty, error) {
	diDB := relationDB.NewDeviceInfoRepo(l.ctx)
	di, err := diDB.FindOneByFilter(ctxs.WithRoot(l.ctx), relationDB.DeviceFilter{
		ProductID:   in.ProductID,
		DeviceNames: []string{in.DeviceName},
	})
	if err != nil {
		return nil, err
	}
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	pi, err := l.svcCtx.ProjectM.ProjectInfoRead(l.ctx, &sys.ProjectWithID{ProjectID: int64(di.ProjectID)})
	if err != nil {
		return nil, err
	}
	//如果是超管有全部权限
	if !uc.AllTenant && (di.TenantCode != di.TenantCode || pi.AdminUserID != uc.UserID || int64(di.ProjectID) != uc.ProjectID) {
		return nil, errors.Permissions
	}
	dpi, err := l.svcCtx.TenantCache.GetData(l.ctx, def.TenantCodeDefault)
	if err != nil {
		return nil, err
	}
	di.TenantCode = def.TenantCodeDefault
	di.ProjectID = stores.ProjectID(dpi.DefaultProjectID)
	if di.ProjectID == 0 {
		di.ProjectID = def.NotClassified
	}
	di.UserID = def.RootNode
	di.AreaID = stores.AreaID(def.NotClassified)
	di.AreaIDPath = def.NotClassifiedPath
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {

		err := relationDB.NewDeviceInfoRepo(tx).Update(ctxs.WithRoot(l.ctx), di)
		if err != nil {
			return err
		}
		err = relationDB.NewUserDeviceShareRepo(tx).DeleteByFilter(l.ctx, relationDB.UserDeviceShareFilter{
			ProductID:  di.ProductID,
			DeviceName: di.DeviceName,
		})
		if err != nil {
			return err
		}
		err = relationDB.NewDeviceProfileRepo(tx).DeleteByFilter(ctxs.WithRoot(l.ctx),
			relationDB.DeviceProfileFilter{Device: devices.Core{
				ProductID:  di.ProductID,
				DeviceName: di.DeviceName,
			}})
		return err
	})
	if err != nil {
		return nil, err
	}
	l.svcCtx.DeviceCache.SetData(l.ctx, devices.Core{
		ProductID:  di.ProductID,
		DeviceName: di.DeviceName,
	}, nil)

	return &dm.Empty{}, err
}
