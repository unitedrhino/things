package devicemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/stores"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"gorm.io/gorm"
	"time"

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
	if !uc.IsAdmin && (di.TenantCode != di.TenantCode || pi.AdminUserID != uc.UserID || int64(di.ProjectID) != uc.ProjectID) {
		return nil, errors.Permissions
	}
	//dpi, err := l.svcCtx.TenantCache.GetData(l.ctx, def.TenantCodeDefault)
	//if err != nil {
	//	return nil, err
	//}
	di.TenantCode = def.TenantCodeDefault
	di.ProjectID = def.NotClassified
	di.UserID = def.RootNode
	di.AreaID = stores.AreaID(def.NotClassified)
	di.AreaIDPath = def.NotClassifiedPath
	if di.FirstBind.Valid && di.FirstBind.Time.After(time.Now().AddDate(0, 0, -1)) { //绑定一天内的不算绑定时间
		di.FirstBind.Valid = false
		di.ExpTime.Valid = false
	}
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
		if err != nil {
			return err
		}
		err = relationDB.NewUserDeviceCollectRepo(tx).DeleteByFilter(l.ctx, relationDB.UserDeviceCollectFilter{Cores: []*devices.Core{
			{
				ProductID:  di.ProductID,
				DeviceName: di.DeviceName,
			},
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
	err = DeleteDeviceTimeData(l.ctx, l.svcCtx, in.ProductID, in.DeviceName)
	err = l.svcCtx.FastEvent.Publish(l.ctx, eventBus.DmDeviceInfoUnbind, &devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName})
	if err != nil {
		l.Error(err)
	}
	return &dm.Empty{}, err
}
