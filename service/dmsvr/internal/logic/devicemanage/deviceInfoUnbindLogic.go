package devicemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

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
	di.AreaID = stores.AreaID(def.NotClassified)
	err = diDB.Update(ctxs.WithRoot(l.ctx), di)

	return &dm.Empty{}, err
}
