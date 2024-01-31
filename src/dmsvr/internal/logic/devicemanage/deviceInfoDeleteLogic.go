package devicemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/devices"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/eventBus"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	DiDB *relationDB.DeviceInfoRepo
}

func NewDeviceInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoDeleteLogic {
	return &DeviceInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
	}
}

// 删除设备
func (l *DeviceInfoDeleteLogic) DeviceInfoDelete(in *dm.DeviceInfoDeleteReq) (*dm.Response, error) {
	di, err := l.DiDB.FindOneByFilter(l.ctx, relationDB.DeviceFilter{ProductID: in.ProductID, DeviceNames: []string{in.DeviceName}})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddDetailf("not find device productId=%s deviceName=%s",
				in.ProductID, in.DeviceName)
		}
		l.Errorf("%s.FindOne err=%+v", utils.FuncName(), err)
		return nil, errors.System.AddDetail(err)
	}
	{ //删除时序数据库中的表数据
		schema, err := l.svcCtx.SchemaRepo.GetSchemaModel(l.ctx, in.ProductID)
		if err != nil {
			l.Errorf("%s.GetSchemaModel err=%+v", utils.FuncName(), err)
			return nil, errors.System.AddDetail(err)
		}
		err = l.svcCtx.HubLogRepo.DropDevice(l.ctx, in.ProductID, in.DeviceName)
		if err != nil {
			l.Errorf("%s.DeviceLogRepo.DeleteDevice err=%v", utils.FuncName(), err)
			return nil, err
		}
		err = l.svcCtx.SchemaManaRepo.DeleteDevice(l.ctx, schema, in.ProductID, in.DeviceName)
		if err != nil {
			l.Errorf("%s.SchemaManaRepo.DeleteDevice err=%v", utils.FuncName(), err)
			return nil, err
		}
		err = l.svcCtx.SDKLogRepo.DropDevice(l.ctx, in.ProductID, in.DeviceName)
		if err != nil {
			l.Errorf("%s.SchemaManaRepo.DeleteDevice err=%v", utils.FuncName(), err)
			return nil, err
		}
	}

	err = l.DiDB.Delete(l.ctx, di.ID)
	if err != nil {
		l.Errorf("%s.DeviceInfo.Delete err=%+v", utils.FuncName(), err)
		return nil, errors.System.AddDetail(err)
	}
	l.svcCtx.ServerMsg.Publish(l.ctx, eventBus.DmDeviceInfoDelete, &devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName})
	return &dm.Response{}, nil
}
