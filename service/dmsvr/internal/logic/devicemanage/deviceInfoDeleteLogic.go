package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/topics"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	DiDB *relationDB.DeviceInfoRepo
}

func NewDeviceInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoDeleteLogic {
	ctx = ctxs.WithDefaultRoot(ctx)
	return &DeviceInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
	}
}

// 删除设备
func (l *DeviceInfoDeleteLogic) DeviceInfoDelete(in *dm.DeviceInfoDeleteReq) (*dm.Empty, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	l.ctx = ctxs.WithDefaultAllProject(l.ctx)
	di, err := l.DiDB.FindOneByFilter(l.ctx, relationDB.DeviceFilter{ProductID: in.ProductID, DeviceNames: []string{in.DeviceName}})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddDetailf("not find device productId=%s deviceName=%s",
				in.ProductID, in.DeviceName)
		}
		l.Errorf("%s.FindOne err=%+v", utils.FuncName(), err)
		return nil, errors.System.AddDetail(err)
	}
	//删除时序数据库中的表数据
	err = DeleteDeviceTimeData(l.ctx, l.svcCtx, in.ProductID, in.DeviceName, DeleteModeAll)
	if err != nil {
		return nil, err
	}
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewDeviceInfoRepo(tx).Delete(l.ctx, di.ID)
		if err != nil {
			l.Errorf("%s.DeviceInfo.Delete err=%+v", utils.FuncName(), err)
			return err
		}
		dev := devices.Core{
			ProductID:  di.ProductID,
			DeviceName: di.DeviceName,
		}
		err = relationDB.NewDeviceProfileRepo(tx).DeleteByFilter(l.ctx, relationDB.DeviceProfileFilter{Device: dev})
		if err != nil {
			l.Errorf("%s.NewDeviceProfileRepo.Delete err=%+v", utils.FuncName(), err)
			return err
		}
		err = relationDB.NewUserDeviceShareRepo(tx).DeleteByFilter(l.ctx, relationDB.UserDeviceShareFilter{
			ProductID:  di.ProductID,
			DeviceName: di.DeviceName,
		})
		if err != nil {
			return err
		}
		err = relationDB.NewUserDeviceCollectRepo(tx).DeleteByFilter(l.ctx, relationDB.UserDeviceCollectFilter{Cores: []*devices.Core{
			{ProductID: di.ProductID, DeviceName: di.DeviceName},
		}})
		if err != nil {
			return err
		}
		err = relationDB.NewDeviceSchemaRepo(tx).DeleteByFilter(l.ctx, relationDB.DeviceSchemaFilter{ProductID: di.ProductID, DeviceName: di.DeviceName})
		if err != nil {
			return err
		}
		err = relationDB.NewGatewayDeviceRepo(tx).DeleteDevAll(l.ctx, dev)
		return err
	})
	if err != nil {
		l.Errorf("%s.DeviceInfo.Delete err=%+v", utils.FuncName(), err)
		return nil, err
	}
	err = l.svcCtx.DeviceCache.SetData(l.ctx, devices.Core{
		ProductID:  di.ProductID,
		DeviceName: di.DeviceName,
	}, nil)
	if err != nil {
		l.Error(err)
	}
	err = l.svcCtx.FastEvent.Publish(l.ctx, topics.DmDeviceInfoDelete, &devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName})
	if err != nil {
		l.Error(err)
	}
	logic.FillAreaDeviceCount(l.ctx, l.svcCtx, string(di.AreaIDPath))
	logic.FillProjectDeviceCount(l.ctx, l.svcCtx, int64(di.ProjectID))
	return &dm.Empty{}, nil
}

type DeleteMode int64

const (
	DeleteModeAll   DeleteMode = iota
	DeleteModeThing            //只删除物模型信息
)

func DeleteDeviceTimeData(ctx context.Context, svcCtx *svc.ServiceContext, productID, deviceName string, mode DeleteMode) error {
	schema, err := svcCtx.DeviceSchemaRepo.GetData(ctx, devices.Core{ProductID: productID, DeviceName: deviceName})
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.GetSchemaModel err=%+v", utils.FuncName(), err)
		return errors.System.AddDetail(err)
	}

	err = svcCtx.SchemaManaRepo.DeleteDevice(ctx, schema, productID, deviceName)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.SchemaManaRepo.DeleteDevice err=%v", utils.FuncName(), err)
		return err
	}
	if mode == DeleteModeThing {
		return nil
	}
	err = svcCtx.HubLogRepo.DeleteDevice(ctx, productID, deviceName)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.HubLogRepo.DeleteDevice err=%v", utils.FuncName(), err)
		return err
	}

	err = svcCtx.SDKLogRepo.DeleteDevice(ctx, productID, deviceName)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.SDKLogRepo.DeleteDevice err=%v", utils.FuncName(), err)
		return err
	}
	err = svcCtx.SendRepo.DeleteDevice(ctx, productID, deviceName)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.SendRepo.DeleteDevice err=%v", utils.FuncName(), err)
		return err
	}
	err = svcCtx.StatusRepo.DeleteDevice(ctx, productID, deviceName)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.StatusRepo.DeleteDevice err=%v", utils.FuncName(), err)
		return err
	}
	return nil
}
