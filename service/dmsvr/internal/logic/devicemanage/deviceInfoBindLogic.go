package devicemanagelogic

import (
	"context"
	"database/sql"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/dmsvr/internal/logic"
	"gitee.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"time"

	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

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
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	projectI, err := l.svcCtx.ProjectCache.GetData(l.ctx, uc.ProjectID)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	if uc.ProjectAuth == nil || uc.ProjectAuth[uc.ProjectID] == nil {
		return nil, errors.Permissions.AddMsg("无权限")
	}
	diDB := relationDB.NewDeviceInfoRepo(l.ctx)
	di, err := diDB.FindOneByFilter(ctxs.WithRoot(l.ctx), relationDB.DeviceFilter{
		ProductID:   in.Device.ProductID,
		DeviceNames: []string{in.Device.DeviceName},
	})
	if err != nil && !errors.Cmp(err, errors.NotFind) {
		l.Error(err)
		return nil, err
	}
	if di == nil {
		di, err = relationDB.NewDeviceInfoRepo(l.ctx).FindOneByFilter(ctxs.WithRoot(l.ctx), relationDB.DeviceFilter{
			DeviceNames: []string{in.Device.DeviceName},
		})
		if err != nil {
			return nil, err
		}
	}

	dpi, err := l.svcCtx.TenantCache.GetData(l.ctx, uc.TenantCode)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	//di.ProjectID=1  di.AreaID=2   dpi.ProjectID=0
	if !((di.TenantCode == def.TenantCodeDefault && di.ProjectID < 3) || int64(di.ProjectID) == uc.ProjectID ||
		int64(di.ProjectID) == dpi.DefaultProjectID) { //如果在其他租户下 则已经被绑定 或 在本租户下,但是不在一个项目下也不允许绑定
		//只有归属于default租户和自己租户的才可以
		l.Infof("DeviceCantBound di:%v uc:%v", utils.Fmt(di), utils.Fmt(uc))
		return nil, errors.DeviceCantBound.WithMsg("设备已被其他用户绑定。如需解绑，请按照相关流程操作。")
	}
	if string(di.TenantCode) == uc.TenantCode &&
		int64(di.ProjectID) == uc.ProjectID { //如果已经绑定到自己名下则不允许重复绑定
		return nil, errors.DeviceBound.WithMsg("设备已存在，请返回设备列表查看该设备")
	}
	di.TenantCode = stores.TenantCode(uc.TenantCode)
	di.ProjectID = stores.ProjectID(uc.ProjectID)
	di.UserID = projectI.AdminUserID
	if in.AreaID == 0 {
		in.AreaID = def.NotClassified
	}
	di.AreaID = stores.AreaID(in.AreaID)
	ai, err := l.svcCtx.AreaCache.GetData(l.ctx, in.AreaID)
	if err != nil {
		return nil, err
	}
	oldAreaIDPath := di.AreaIDPath
	di.AreaIDPath = ai.AreaIDPath

	if !di.FirstBind.Valid { //没有绑定过需要绑定
		di.FirstBind = sql.NullTime{Time: time.Now(), Valid: true}
	}
	pi, err := l.svcCtx.ProductCache.GetData(l.ctx, di.ProductID)
	if err != nil && !errors.Cmp(err, errors.NotFind) {
		l.Error(err)
		return nil, err
	}
	if pi.TrialTime.GetValue() != 0 && !di.ExpTime.Valid {
		di.ExpTime = sql.NullTime{
			Time:  time.Now().Add(time.Hour * 24 * time.Duration(pi.TrialTime.GetValue())),
			Valid: true,
		}
	}
	err = diDB.Update(ctxs.WithRoot(l.ctx), di)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	l.svcCtx.DeviceCache.SetData(l.ctx, devices.Core{
		ProductID:  di.ProductID,
		DeviceName: di.DeviceName,
	}, nil)
	{ //清除之前的日志
		err = l.svcCtx.SendRepo.DeleteDevice(l.ctx, di.ProductID, di.DeviceName)
		if err != nil {
			l.Error(err)
		}
		err = l.svcCtx.StatusRepo.DeleteDevice(l.ctx, di.ProductID, di.DeviceName)
		if err != nil {
			l.Error(err)
		}
	}
	logic.FillAreaDeviceCount(l.ctx, l.svcCtx, ai.AreaIDPath, oldAreaIDPath)
	logic.FillProjectDeviceCount(l.ctx, l.svcCtx, int64(di.ProjectID))
	return &dm.Empty{}, err
}
