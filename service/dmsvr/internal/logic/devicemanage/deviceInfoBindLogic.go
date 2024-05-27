package devicemanagelogic

import (
	"context"
	"database/sql"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"time"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoBindLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceInfoBindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoBindLogic {
	return &DeviceInfoBindLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceInfoBindLogic) DeviceInfoBind(in *dm.DeviceInfoBindReq) (*dm.Empty, error) {
	if in.ProtocolCode != "" {
		pi, err := relationDB.NewProductInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.ProductFilter{
			ProtocolConf: map[string]string{"productID": in.Device.ProductID},
			ProtocolCode: in.ProtocolCode,
		})
		if err != nil {
			return nil, err
		}
		in.Device.ProductID = pi.ProductID
	}
	diDB := relationDB.NewDeviceInfoRepo(l.ctx)
	di, err := diDB.FindOneByFilter(ctxs.WithRoot(l.ctx), relationDB.DeviceFilter{
		ProductID:   in.Device.ProductID,
		DeviceNames: []string{in.Device.DeviceName},
	})
	if err != nil {
		return nil, err
	}
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	if (di.TenantCode != def.TenantCodeDefault && string(di.TenantCode) != uc.TenantCode) || (string(di.TenantCode) == uc.TenantCode &&
		int64(di.ProjectID) != uc.ProjectID) { //如果在其他租户下 则已经被绑定 或 在本租户下,但是不在一个项目下也不允许绑定
		//只有归属于default租户和自己租户的才可以
		return nil, errors.DeviceCantBound
	}
	if string(di.TenantCode) == uc.TenantCode &&
		int64(di.ProjectID) == uc.ProjectID { //如果已经绑定到自己名下则不允许重复绑定
		return nil, errors.DeviceBound
	}
	di.TenantCode = stores.TenantCode(uc.TenantCode)
	di.ProjectID = stores.ProjectID(uc.ProjectID)
	if in.AreaID == 0 {
		in.AreaID = def.NotClassified
	}
	di.AreaID = stores.AreaID(in.AreaID)
	if di.AreaID == 0 {
		di.AreaID = def.NotClassified
	}
	if di.FirstBind.Valid {
		di.FirstBind = sql.NullTime{Time: time.Now(), Valid: true}
	}
	err = diDB.Update(ctxs.WithRoot(l.ctx), di)
	l.svcCtx.DeviceCache.SetData(l.ctx, devices.Core{
		ProductID:  di.ProductID,
		DeviceName: di.DeviceName,
	}, nil)
	return &dm.Empty{}, err
}
