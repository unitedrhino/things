package devicemanagelogic

import (
	"context"
	"database/sql"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/domain/deviceMsg/msgOta"
	"gitee.com/i-Things/share/domain/deviceMsg/msgSdkLog"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/dmExport"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
	DiDB *relationDB.DeviceInfoRepo
}

func NewDeviceInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoUpdateLogic {
	return &DeviceInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
	}
}

func (l *DeviceInfoUpdateLogic) SetDevicePoByDto(old *relationDB.DmDeviceInfo, data *dm.DeviceInfo) error {
	if data.AreaID != 0 && data.AreaID != int64(old.AreaID) {
		old.AreaID = stores.AreaID(data.AreaID)
		err := l.svcCtx.StatusRepo.ModifyDeviceArea(l.ctx, devices.Core{
			ProductID:  data.ProductID,
			DeviceName: data.DeviceName,
		}, data.AreaID)
		if err != nil {
			l.Error(err)
			return errors.Database.AddDetail(err)
		}
		err = l.svcCtx.SendRepo.ModifyDeviceArea(l.ctx, devices.Core{
			ProductID:  data.ProductID,
			DeviceName: data.DeviceName,
		}, data.AreaID)
		if err != nil {
			l.Error(err)
			return errors.Database.AddDetail(err)
		}
	}
	if data.ProjectID != 0 && data.ProjectID != int64(old.ProjectID) {
		err := l.svcCtx.StatusRepo.ModifyDeviceProject(l.ctx, devices.Core{
			ProductID:  data.ProductID,
			DeviceName: data.DeviceName,
		}, data.ProjectID)
		if err != nil {
			l.Error(err)
			return errors.Database.AddDetail(err)
		}
		err = l.svcCtx.SendRepo.ModifyDeviceProject(l.ctx, devices.Core{
			ProductID:  data.ProductID,
			DeviceName: data.DeviceName,
		}, data.ProjectID)
		if err != nil {
			l.Error(err)
			return errors.Database.AddDetail(err)
		}
		old.ProjectID = stores.ProjectID(data.ProjectID)
	}

	if data.Tags != nil {
		old.Tags = data.Tags
	}
	if data.ProtocolConf != nil {
		old.ProtocolConf = data.ProtocolConf
	}
	if data.SchemaAlias != nil {
		old.SchemaAlias = data.SchemaAlias
	}
	if data.LogLevel != def.Unknown {
		old.LogLevel = data.LogLevel
	}

	if data.Imei != "" {
		old.Imei = data.Imei
	}
	if data.Mac != "" {
		old.Mac = data.Mac
	}
	if data.Version != nil && old.Version != data.Version.GetValue() {
		//如果不一样则需要判断是否是ota升级的,如果是,则需要更新升级状态
		df, err := relationDB.NewOtaFirmwareDeviceRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.OtaFirmwareDeviceFilter{
			ProductID:   old.ProductID,
			DeviceNames: []string{old.DeviceName},
			Statues:     []int64{msgOta.DeviceStatusInProgress, msgOta.DeviceStatusNotified},
		})
		if err != nil {
			if !errors.Cmp(err, errors.NotFind) {
				return err
			}
		} else {
			if df.DestVersion == data.Version.GetValue() { //版本号一致才是升级成功
				old.Version = data.Version.GetValue()
				df.Step = 100
				df.Status = msgOta.DeviceStatusSuccess
				err := relationDB.NewOtaFirmwareDeviceRepo(l.ctx).Update(l.ctx, df)
				if err != nil {
					return err
				}
			}
		}
	}
	if data.HardInfo != "" {
		old.HardInfo = data.HardInfo
	}
	if data.SoftInfo != "" {
		old.SoftInfo = data.SoftInfo
	}

	if data.Rssi != nil {
		old.Rssi = data.Rssi.GetValue()
	}

	if data.IsOnline != def.Unknown {
		old.IsOnline = data.IsOnline
		if data.IsOnline == def.True { //需要处理第一次上线的情况,一般在网关代理登录时需要处理
			now := sql.NullTime{
				Valid: true,
				Time:  time.Now(),
			}
			if old.FirstLogin.Valid == false {
				old.FirstLogin = now
			}
			old.LastLogin = now
		}
	}

	if data.Address != nil {
		old.Address = data.Address.Value
	}
	if data.Position != nil {
		old.Position = logic.ToStorePoint(data.Position)
	}

	if data.DeviceAlias != nil {
		old.DeviceAlias = data.DeviceAlias.Value
	}
	if data.MobileOperator != 0 {
		old.MobileOperator = data.MobileOperator
	}
	if data.Iccid != nil {
		old.Iccid = utils.AnyToNullString(data.Iccid)
	}
	if data.Phone != nil {
		old.Phone = utils.AnyToNullString(data.Phone)
	}
	return nil
}

// 更新设备
func (l *DeviceInfoUpdateLogic) DeviceInfoUpdate(in *dm.DeviceInfo) (*dm.Empty, error) {
	if in.ProductID == "" && in.ProductName != "" { //通过唯一的产品名 查找唯一的产品ID
		if pid, err := l.PiDB.FindOneByFilter(l.ctx, relationDB.ProductFilter{ProductNames: []string{in.ProductName}}); err != nil {
			return nil, err
		} else {
			in.ProductID = pid.ProductID
		}
	}
	dmDiPo, err := l.DiDB.FindOneByFilter(l.ctx, relationDB.DeviceFilter{ProductID: in.ProductID, DeviceNames: []string{in.DeviceName}})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.NotFind.AddDetailf("not find device productID=%s deviceName=%s",
				in.ProductID, in.DeviceName)
		}
		return nil, errors.Database.AddDetail(err)
	}

	l.SetDevicePoByDto(dmDiPo, in)

	err = l.DiDB.Update(l.ctx, dmDiPo)
	if err != nil {
		l.Errorf("DeviceInfoUpdate.DeviceInfo.Update err=%+v", err)
		return nil, err
	}
	err = l.svcCtx.DeviceCache.SetData(l.ctx, dmExport.GenDeviceInfoKey(dmDiPo.ProductID, dmDiPo.DeviceName), logic.ToDeviceInfo(dmDiPo))
	if err != nil {
		l.Error(err)
	}
	if in.LogLevel != def.Unknown {
		di, err := l.DiDB.FindOneByFilter(l.ctx, relationDB.DeviceFilter{ProductID: in.ProductID, DeviceNames: []string{in.DeviceName}, WithProduct: true})
		if err != nil {
			return nil, err
		}
		resp := deviceMsg.NewRespCommonMsg(l.ctx, deviceMsg.GetStatus, "")
		resp.Data = map[string]any{"logLevel": di.LogLevel}

		msg := deviceMsg.PublishMsg{
			Handle:     devices.Log,
			Type:       msgSdkLog.TypeUpdate,
			Payload:    resp.AddStatus(errors.OK).Bytes(),
			Timestamp:  time.Now().UnixMilli(),
			ProductID:  di.ProductID,
			DeviceName: di.DeviceName,
		}
		er := l.svcCtx.PubDev.PublishToDev(l.ctx, &msg)
		if er != nil {
			l.Errorf("%s.PublishToDev failure err:%v", utils.FuncName(), er)
			return nil, err
		}
	}
	return &dm.Empty{}, nil
}
