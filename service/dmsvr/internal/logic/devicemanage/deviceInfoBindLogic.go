package devicemanagelogic

import (
	"context"
	"database/sql"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"time"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

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
	pi, err := l.svcCtx.ProductCache.GetData(l.ctx, in.Device.ProductID)
	if err != nil {
		l.Error(err)
		return nil, err
	}
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
			if !errors.Cmp(err, errors.NotFind) {
				return nil, err
			}

			if !(pi.NetType == def.NetBle && pi.AutoRegister == def.AutoRegAuto) {
				return nil, errors.NotFind
			}
			//如果是蓝牙模式并且打开了自动注册,那么绑定的时候需要创建该设备
			_, err = NewDeviceInfoCreateLogic(ctxs.WithProjectID(ctxs.WithAdmin(l.ctx), def.NotClassified), l.svcCtx).DeviceInfoCreate(&dm.DeviceInfo{ProductID: in.Device.ProductID, DeviceName: in.Device.DeviceName})
			if err != nil {
				return nil, err
			}
			di, err = relationDB.NewDeviceInfoRepo(l.ctx).FindOneByFilter(ctxs.WithRoot(l.ctx), relationDB.DeviceFilter{
				ProductID:   in.Device.ProductID,
				DeviceNames: []string{in.Device.DeviceName},
			})
			if err != nil {
				return nil, err
			}
		}
	}

	dpi, err := l.svcCtx.TenantCache.GetData(l.ctx, def.TenantCodeDefault)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	//di.ProjectID=1  di.AreaID=2   dpi.ProjectID=0
	if !((di.TenantCode == def.TenantCodeDefault && di.ProjectID < 3) || int64(di.ProjectID) == uc.ProjectID ||
		int64(di.ProjectID) == dpi.DefaultProjectID) && !(pi.IsCanCoverBindDevice == def.True) { //如果在其他租户下 则已经被绑定 或 在本租户下,但是不在一个项目下也不允许绑定
		//只有归属于default租户和自己租户的才可以
		l.Infof("DeviceCantBound di:%v uc:%v", utils.Fmt(di), utils.Fmt(uc))
		return nil, errors.DeviceCantBound.WithMsg("设备已被其他用户绑定。如需解绑，请按照相关流程操作。")
	}
	if string(di.TenantCode) == uc.TenantCode &&
		int64(di.ProjectID) == uc.ProjectID { //如果已经绑定到自己名下
		return &dm.Empty{}, err
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
	di.AreaIDPath = stores.AreaIDPath(ai.AreaIDPath)

	if !di.FirstBind.Valid { //没有绑定过需要绑定
		di.FirstBind = sql.NullTime{Time: time.Now(), Valid: true}
	}
	di.LastBind = sql.NullTime{Time: time.Now(), Valid: true}
	pc, err := l.svcCtx.ProductCache.GetData(l.ctx, di.ProductID)
	if err != nil && !errors.Cmp(err, errors.NotFind) {
		l.Error(err)
		return nil, err
	}
	if pc.TrialTime.GetValue() != 0 && !di.ExpTime.Valid {
		di.ExpTime = sql.NullTime{
			Time:  time.Now().Add(time.Hour * 24 * time.Duration(pc.TrialTime.GetValue())),
			Valid: true,
		}
	}
	if pc.NetType == def.NetBle { //蓝牙绑定了就是上线
		di.IsOnline = def.True
		di.Status = def.DeviceStatusOnline
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
	logic.FillAreaDeviceCount(l.ctx, l.svcCtx, ai.AreaIDPath, string(oldAreaIDPath))
	logic.FillProjectDeviceCount(l.ctx, l.svcCtx, int64(di.ProjectID))
	return &dm.Empty{}, err
}
