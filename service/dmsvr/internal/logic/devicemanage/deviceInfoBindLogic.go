package devicemanagelogic

import (
	"context"
	"database/sql"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/product"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/topics"
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

func isAllowedChar(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_' ||
		c == '-'
}

func filterAllowedChars(input string) string {
	var result []rune
	for _, c := range input {
		if isAllowedChar(c) {
			result = append(result, c)
		}
	}
	return string(result)
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
			DeviceNames: []string{in.Device.DeviceName, filterAllowedChars(in.Device.DeviceName)}, //兼容打印错误
		})
		if err != nil {
			if !errors.Cmp(err, errors.NotFind) {
				return nil, err
			}
			if !(pi.NetType == def.NetBle && pi.AutoRegister == def.AutoRegAuto) {
				return nil, errors.NotFind
			}
			//如果是蓝牙模式并且打开了自动注册,那么绑定的时候需要创建该设备
			_, err = NewDeviceInfoCreateLogic(ctxs.WithProjectID(ctxs.WithAdmin(l.ctx), def.NotClassified), l.svcCtx).
				DeviceInfoCreate(&dm.DeviceInfo{ProductID: in.Device.ProductID, DeviceName: in.Device.DeviceName, IsOnline: def.True})
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
	if pi.BindLevel < 3 && di.IsOnline != def.True { //如果是中绑定和强绑定,如果设备不在线,不允许绑定
		return nil, errors.NotOnline
	}
	if string(di.TenantCode) == uc.TenantCode &&
		int64(di.ProjectID) == uc.ProjectID { //如果已经绑定到自己名下则不允许重复绑定
		if pi.BindLevel == product.BindLeveHard1 {
			return nil, errors.DeviceBound.WithMsg("设备已存在，请返回设备列表查看该设备")
		} else {
			return &dm.Empty{}, nil
		}
	}

	dpi, err := l.svcCtx.TenantCache.GetData(l.ctx, def.TenantCodeDefault)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	//是否可以绑定校验
	if pi.BindLevel == product.BindLeveMiddle2 && in.Token != "" {
		tk, err := l.svcCtx.DeviceBindToken.GetData(l.ctx, in.Token)
		if err != nil {
			return nil, errors.NotFind.AddMsg("未发现配网token").AddDetail(err)
		}
		if tk.UserID != uc.UserID {
			return nil, errors.Permissions.AddMsg("配网和绑定的用户不一致")
		}
	} else {
		if !((di.TenantCode == def.TenantCodeDefault && di.ProjectID < 3) || int64(di.ProjectID) == uc.ProjectID ||
			int64(di.ProjectID) == dpi.DefaultProjectID) { //如果在其他租户下 则已经被绑定 或 在本租户下,但是不在一个项目下也不允许绑定
			//只有归属于default租户和自己租户的才可以
			l.Infof("DeviceCantBound di:%v uc:%v", utils.Fmt(di), utils.Fmt(uc))
			return nil, errors.DeviceCantBound.WithMsg("设备已被其他用户绑定。如需解绑，请按照相关流程操作。")
		}
	}

	di.TenantCode = dataType.TenantCode(uc.TenantCode)
	di.ProjectID = dataType.ProjectID(uc.ProjectID)
	di.UserID = projectI.AdminUserID
	if in.AreaID == 0 {
		in.AreaID = def.NotClassified
	}
	di.AreaID = dataType.AreaID(in.AreaID)
	ai, err := l.svcCtx.AreaCache.GetData(l.ctx, in.AreaID)
	if err != nil {
		return nil, err
	}
	oldAreaIDPath := di.AreaIDPath
	di.AreaIDPath = dataType.AreaIDPath(ai.AreaIDPath)

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
	err = l.svcCtx.AbnormalRepo.UpdateDevice(l.ctx, []*devices.Core{
		{ProductID: di.ProductID, DeviceName: di.DeviceName}}, devices.Affiliation{TenantCode: string(di.TenantCode),
		ProjectID: int64(di.ProjectID), AreaID: int64(di.AreaID), AreaIDPath: string(di.AreaIDPath)})
	if err != nil {
		l.Error(err)
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
	dev := devices.Core{ProductID: di.ProductID, DeviceName: di.DeviceName}
	er := l.svcCtx.FastEvent.Publish(l.ctx, topics.DmDeviceInfoUpdate, &dev)
	if er != nil {
		l.Error(er)
	}
	BindChange(l.ctx, l.svcCtx, pi, dev, int64(di.ProjectID))
	return &dm.Empty{}, err
}
