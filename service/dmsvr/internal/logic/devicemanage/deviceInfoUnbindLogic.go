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
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gorm.io/gorm"
	"time"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

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
	if err != nil && !errors.Cmp(err, errors.NotFind) { //解绑的时候家庭已经不存在了也需要能正确解绑
		return nil, err
	}
	adminUserID := di.UserID
	if pi != nil {
		adminUserID = pi.AdminUserID
	}
	//如果是超管有全部权限
	if !uc.IsAdmin && (di.TenantCode != di.TenantCode || adminUserID != uc.UserID || int64(di.ProjectID) != uc.ProjectID) {
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
		pc, err := l.svcCtx.ProductCache.GetData(l.ctx, di.ProductID)
		if err != nil {
			return nil, err
		}
		if pc.TrialTime != nil && di.ExpTime.Valid { //如果设备的有效期大于从当前算起的有效期,那说明充值过,这时候不能清除过期时间
			expTime := time.Now().Add(time.Hour * 24 * time.Duration(pc.TrialTime.GetValue()))
			if expTime.After(di.ExpTime.Time) {
				di.FirstBind.Valid = false
				di.ExpTime.Valid = false
			}
		}
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
	if di.DeviceType == def.DeviceTypeGateway { //网关类型的需要解绑子设备
		ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
			subs, err := relationDB.NewGatewayDeviceRepo(ctx).FindByFilter(l.ctx, relationDB.GatewayDeviceFilter{Gateway: &devices.Core{
				ProductID:  di.ProductID,
				DeviceName: di.DeviceName,
			}}, nil)
			if err != nil {
				logx.WithContext(ctx).Error(err)
				return
			}
			for _, sub := range subs {
				_, err := NewDeviceInfoUnbindLogic(ctx, l.svcCtx).DeviceInfoUnbind(&dm.DeviceCore{
					ProductID:  sub.ProductID,
					DeviceName: sub.DeviceName,
				})
				if err != nil {
					logx.WithContext(ctx).Error(err)
					continue
				}
			}
		})
	}
	return &dm.Empty{}, err
}
