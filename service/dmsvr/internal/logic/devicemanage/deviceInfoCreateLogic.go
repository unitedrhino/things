package devicemanagelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/atomic"
)

type DeviceInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
	DiDB *relationDB.DeviceInfoRepo
}

var randID atomic.Uint32

func GenID() uint32 {
	return randID.Inc() % 100
}

func NewDeviceInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoCreateLogic {
	return &DeviceInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
	}
}

/*
发现返回true 没有返回false
*/
func (l *DeviceInfoCreateLogic) CheckDevice(in *dm.DeviceInfo) (bool, error) {
	_, err := relationDB.NewDeviceInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.DeviceFilter{ProductID: in.ProductID, DeviceNames: []string{in.DeviceName}})
	if err == nil {
		return true, nil
	}
	if errors.Cmp(err, errors.NotFind) {
		return false, nil
	}
	return false, err
}

/*
发现返回true 没有返回false
*/
func (l *DeviceInfoCreateLogic) CheckProduct(in *dm.DeviceInfo) (*dm.ProductInfo, error) {
	pi, err := l.svcCtx.ProductCache.GetData(l.ctx, in.ProductID)
	if err == nil {
		return pi, nil
	}
	if errors.Cmp(err, errors.NotFind) {
		return nil, nil
	}
	return nil, err
}

// 新增设备
func (l *DeviceInfoCreateLogic) DeviceInfoCreate(in *dm.DeviceInfo) (resp *dm.Empty, err error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	l.ctx = ctxs.WithDefaultAllProject(l.ctx)
	if in.ProductID == "" && in.ProductName != "" { //通过唯一的产品名 查找唯一的产品ID
		if pid, err := l.PiDB.FindOneByFilter(l.ctx, relationDB.ProductFilter{ProductNames: []string{in.ProductName}}); err != nil {
			return nil, err
		} else {
			in.ProductID = pid.ProductID
		}
	}

	find, err := l.CheckDevice(in)
	if err != nil {
		l.Errorf("%s.CheckDevice in=%v\n", utils.FuncName(), in)
		return nil, errors.Database.AddDetail(err)
	} else if find == true {
		return nil, errors.Duplicate.WithMsgf("设备名称重复:%s", in.DeviceName).AddDetail("DeviceName:" + in.DeviceName)
	}

	pi, err := l.CheckProduct(in)
	if err != nil {
		l.Errorf("%s.CheckProduct in=%v", utils.FuncName(), in)
		return nil, err
	} else if pi == nil {
		return nil, errors.Parameter.AddDetail("not find product id:" + cast.ToString(in.ProductID))
	}
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	projectID := stores.ProjectID(uc.ProjectID)
	areaID := stores.AreaID(def.NotClassified)
	if projectID == 0 || projectID == def.NotClassified { //如果没有传项目,则分配到未分类项目中
		ti, err := l.svcCtx.TenantCache.GetData(l.ctx, uc.TenantCode)
		if err != nil {
			return nil, err
		}
		projectID = stores.ProjectID(ti.DefaultProjectID)
		if ti.DefaultAreaID != 0 {
			areaID = stores.AreaID(ti.DefaultAreaID)
		}
	}

	ai, err := l.svcCtx.AreaCache.GetData(l.ctx, int64(areaID))
	if err != nil {
		return nil, err
	}
	areaIDPath := ai.AreaIDPath

	di := relationDB.DmDeviceInfo{
		ProjectID:   projectID,
		ProductID:   in.ProductID,  // 产品id
		DeviceName:  in.DeviceName, // 设备名称
		Position:    logic.ToStorePoint(in.Position),
		AreaID:      areaID, //设备默认都是未分类
		AreaIDPath:  areaIDPath,
		Status:      def.DeviceStatusInactive,
		IsEnable:    def.True,
		RatedPower:  in.RatedPower,
		Distributor: utils.Copy2[stores.IDPathWithUpdate](in.Distributor),
		UserID:      def.RootNode,
	}
	if di.Distributor.ID == 0 {
		di.Distributor.ID = def.RootNode
	}
	if in.IsEnable != 0 {
		di.IsEnable = in.IsEnable
	}
	if in.Secret != "" {
		di.Secret = in.Secret
	} else {
		di.Secret = utils.GetRandomBase64(20)
	}

	di.Tags = in.Tags
	if di.Tags == nil {
		di.Tags = map[string]string{}
	}
	di.ProtocolConf = in.ProtocolConf
	if di.ProtocolConf == nil {
		di.ProtocolConf = map[string]string{}
	}
	di.SchemaAlias = in.SchemaAlias
	if di.SchemaAlias == nil {
		di.SchemaAlias = map[string]string{}
	}

	if in.Rssi != nil {
		di.Rssi = in.Rssi.GetValue()
	}

	if in.LogLevel != def.Unknown {
		di.LogLevel = def.LogClose
	}

	if in.Address != nil {
		di.Address = in.Address.Value
	}

	if in.DeviceAlias != nil {
		di.DeviceAlias = in.DeviceAlias.Value
	} else {
		di.DeviceAlias = fmt.Sprintf("%s%d", pi.ProductName, GenID())
	}

	if in.MobileOperator != 0 {
		di.MobileOperator = in.MobileOperator
	}

	if in.Iccid != nil {
		di.Iccid = utils.AnyToNullString(in.Iccid)
	}

	if in.Phone != nil {
		di.Phone = utils.AnyToNullString(in.Phone)
	}

	if ctxs.IsRoot(l.ctx) == nil {
		if in.Status != 0 {
			di.Status = in.Status
		}
		if in.IsOnline != 0 {
			di.IsOnline = in.IsOnline
		}
	}

	err = l.InitDevice(devices.Info{
		ProductID:  di.ProductID,
		DeviceName: di.DeviceName,
		TenantCode: string(di.TenantCode),
		ProjectID:  int64(di.ProjectID),
		AreaID:     int64(di.AreaID),
	})
	if err != nil {
		return nil, err
	}
	err = l.DiDB.Insert(l.ctx, &di)
	if err != nil {
		l.Errorf("AddDevice.DeviceInfo.Insert err=%+v", err)
		return nil, errors.Database.AddDetail(err)
	}
	logic.FillAreaDeviceCount(l.ctx, l.svcCtx, areaIDPath)
	return &dm.Empty{}, nil
}

func (l *DeviceInfoCreateLogic) InitDevice(in devices.Info) error {
	if in.TenantCode == "" {
		in.TenantCode = ctxs.GetUserCtxNoNil(l.ctx).TenantCode
	}
	pt, err := l.svcCtx.SchemaRepo.GetData(l.ctx, in.ProductID)
	if err != nil {
		return errors.System.AddDetail(err)
	}
	err = l.svcCtx.SchemaManaRepo.InitDevice(l.ctx, pt, in.ProductID, in.DeviceName)
	if err != nil {
		return errors.Database.AddDetail(err)
	}
	err = l.svcCtx.SDKLogRepo.InitDevice(l.ctx, in)
	if err != nil {
		return errors.Database.AddDetail(err)
	}
	err = l.svcCtx.StatusRepo.InitDevice(l.ctx, in)
	if err != nil {
		return errors.Database.AddDetail(err)
	}
	err = l.svcCtx.SendRepo.InitDevice(l.ctx, in)
	if err != nil {
		return errors.Database.AddDetail(err)
	}
	err = l.svcCtx.DeviceCache.SetData(l.ctx, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}, nil)
	if err != nil {
		l.Error(err)
	}
	return nil
}
