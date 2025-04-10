package devicemanagelogic

import (
	"context"
	"database/sql"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	otamanagelogic "gitee.com/unitedrhino/things/service/dmsvr/internal/logic/otamanage"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/application"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgOta"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgSdkLog"
	"gitee.com/unitedrhino/things/share/topics"
	"gitee.com/unitedrhino/things/share/userSubscribe"
	"github.com/spf13/cast"
	"sync"
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
	uc := ctxs.GetUserCtx(l.ctx)
	var isUpdateTag bool
	if data.AreaID != 0 && data.AreaID != int64(old.AreaID) {
		old.AreaID = dataType.AreaID(data.AreaID)
		ai, err := l.svcCtx.AreaCache.GetData(l.ctx, data.AreaID)
		if err != nil {
			return err
		}
		ctxs.GoNewCtx(l.ctx, func(ctx2 context.Context) {
			time.Sleep(2 * time.Second)
			logic.FillAreaDeviceCount(l.ctx, l.svcCtx, ai.AreaIDPath, string(old.AreaIDPath))
		})
		old.AreaIDPath = dataType.AreaIDPath(ai.AreaIDPath)
		isUpdateTag = true
	}
	if data.ProjectID != 0 && data.ProjectID != int64(old.ProjectID) {
		ctxs.GoNewCtx(l.ctx, func(ctx2 context.Context) {
			time.Sleep(2 * time.Second)
			logic.FillProjectDeviceCount(ctx2, l.svcCtx, data.ProjectID, int64(old.ProjectID))
		})
		old.ProjectID = dataType.ProjectID(data.ProjectID)
		isUpdateTag = true
	}
	if data.Sort != 0 {
		old.Sort = data.Sort
	}
	if uc.IsAdmin && data.ExpTime != nil {
		if data.ExpTime.Value == 0 && old.UserID > def.RootNode { //如果是主动传0且被绑定了则需要看下过期时间是否正确
			pi, err := l.svcCtx.ProductCache.GetData(l.ctx, data.ProductID)
			if err != nil && !errors.Cmp(err, errors.NotFind) {
				l.Error(err)
				return err
			}
			if pi.TrialTime.GetValue() != 0 {
				old.ExpTime = sql.NullTime{
					Time:  time.Now().Add(time.Hour * 24 * time.Duration(pi.TrialTime.GetValue())),
					Valid: true,
				}
			}
		} else {
			old.ExpTime = utils.ToNullTime2(data.ExpTime)
		}
		if old.Status == def.DeviceStatusArrearage && (!old.ExpTime.Valid || old.ExpTime.Time.After(time.Now())) {
			//如果设备是欠费状态需要修改为正常状态
			old.Status = old.IsOnline + 1
		}
		if old.Status != def.DeviceStatusArrearage && old.ExpTime.Valid && old.ExpTime.Time.Before(time.Now()) {
			old.Status = def.DeviceStatusArrearage
		}
	}
	if data.Desc != nil {
		old.Desc = data.Desc.GetValue()
	}
	if data.Tags != nil {
		old.Tags = data.Tags
	}
	if uc.IsAdmin && data.ProtocolConf != nil {
		old.ProtocolConf = data.ProtocolConf
	}
	if uc.IsAdmin && data.SubProtocolConf != nil {
		old.SubProtocolConf = data.SubProtocolConf
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
	if uc.IsAdmin && data.Version != nil {
		if old.Version != data.Version.GetValue() {
			old.Version = data.Version.GetValue()
			//如果不一样则需要判断是否是ota升级的,如果是,则需要更新升级状态
			dfs, err := relationDB.NewOtaFirmwareDeviceRepo(l.ctx).FindByFilter(l.ctx, relationDB.OtaFirmwareDeviceFilter{
				ProductID:   old.ProductID,
				DeviceNames: []string{old.DeviceName},
				DestVersion: data.Version.GetValue(),
				Statues: []int64{msgOta.DeviceStatusConfirm, msgOta.DeviceStatusQueued,
					msgOta.DeviceStatusNotified, msgOta.DeviceStatusInProgress},
			}, nil)
			if err != nil {
				if !errors.Cmp(err, errors.NotFind) {
					return err
				}
			} else {
				var once sync.Once
				for _, df := range dfs {
					if df.DestVersion == data.Version.GetValue() { //版本号一致才是升级成功
						df.Step = 100
						df.Status = msgOta.DeviceStatusSuccess
						df.Detail = "升级成功"
						err := relationDB.NewOtaFirmwareDeviceRepo(l.ctx).Update(l.ctx, df)
						if err != nil {
							return err
						}
						old.NeedConfirmVersion = ""
						old.NeedConfirmJobID = 0
						once.Do(func() {
							appMsg := application.OtaReport{
								Device:    devices.Core{ProductID: old.ProductID, DeviceName: old.DeviceName},
								Timestamp: time.Now().UnixMilli(), Status: df.Status, Detail: df.Detail, Step: df.Step,
							}
							err = l.svcCtx.UserSubscribe.Publish(l.ctx, userSubscribe.DeviceOtaReport, appMsg, map[string]any{
								"productID":  old.ProductID,
								"deviceName": old.DeviceName,
							}, map[string]any{
								"projectID": old.ProjectID,
							}, map[string]any{
								"projectID": cast.ToString(old.ProjectID),
								"areaID":    cast.ToString(old.AreaID),
							})
						})
					}
				}

			}
		} else {
			//检查是否有待升级的任务并推送
			dfs, err := relationDB.NewOtaFirmwareDeviceRepo(l.ctx).FindByFilter(l.ctx, relationDB.OtaFirmwareDeviceFilter{
				ProductID:    old.ProductID,
				DeviceNames:  []string{old.DeviceName},
				DestVersion:  data.Version.GetValue(),
				WithJob:      true,
				WithFirmware: true,
				Statues:      []int64{msgOta.DeviceStatusQueued},
			}, nil)
			if err != nil {
				if !errors.Cmp(err, errors.NotFind) {
					return err
				}
			}
			if len(dfs) > 0 && dfs[0].Job != nil && dfs[0].Firmware != nil {
				dfs[0].Job.Firmware = dfs[0].Firmware
				payload, err := otamanagelogic.GenPayload(l.ctx, l.svcCtx, dfs[0].Job)
				if err != nil {
					return err
				}
				pi, err := l.svcCtx.ProductCache.GetData(l.ctx, dfs[0].Firmware.ProductID)
				if err != nil {
					return err
				}
				err = otamanagelogic.PushMessageToDevice(l.ctx, l.svcCtx, dfs[0], pi, payload)
				if err != nil {
					l.Error(err)
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
	if data.RatedPower != 0 {
		old.RatedPower = data.RatedPower
	}
	if uc.IsAdmin && data.IsOnline != def.Unknown {
		old.IsOnline = data.IsOnline
		if old.Status <= def.DeviceStatusOffline {
			old.Status = data.IsOnline + 1
		}
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
	if uc.IsAdmin && data.Status != 0 {
		old.Status = data.Status
	}

	if data.Address != nil {
		old.Address = data.Address.Value
	}
	if data.Adcode != nil {
		old.Adcode = data.Adcode.Value
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
	if uc.IsAdmin && data.Distributor != nil {
		old.Distributor = utils.Copy2[stores.IDPathWithUpdate](data.Distributor)
		old.Distributor.UpdatedTime = time.Now()
	}
	if uc.IsAdmin && data.Dept != nil {
		old.DeptID = dataType.DeptID(data.Dept.Id)
		old.DeptIDPath = dataType.DeptIDPath(data.Dept.IdPath)
		old.DeptUpdatedTime = time.Now()
	}
	if data.Iccid != nil {
		old.Iccid = utils.AnyToNullString(data.Iccid)
	}
	if data.Phone != nil {
		old.Phone = utils.AnyToNullString(data.Phone)
	}

	if isUpdateTag {
		err := l.svcCtx.AbnormalRepo.UpdateDevice(l.ctx, []*devices.Core{
			{ProductID: old.ProductID, DeviceName: old.DeviceName}},
			devices.Affiliation{ProjectID: int64(old.ProjectID), AreaID: int64(old.AreaID), AreaIDPath: string(old.AreaIDPath)})
		if err != nil {
			return err
		}
	}
	return nil
}

// 更新设备
func (l *DeviceInfoUpdateLogic) DeviceInfoUpdate(in *dm.DeviceInfo) (*dm.Empty, error) {
	l.ctx = ctxs.WithDefaultAllProject(l.ctx)
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

	err = l.SetDevicePoByDto(dmDiPo, in)
	if err != nil {
		return nil, err
	}

	err = l.DiDB.Update(l.ctx, dmDiPo)
	if err != nil {
		l.Errorf("DeviceInfoUpdate.DeviceInfo.Update err=%+v", err)
		return nil, err
	}
	err = l.svcCtx.DeviceCache.SetData(l.ctx, devices.Core{
		ProductID:  dmDiPo.ProductID,
		DeviceName: dmDiPo.DeviceName,
	}, nil)
	if err != nil {
		l.Error(err)
	}
	if in.LogLevel != def.Unknown {
		di, err := l.DiDB.FindOneByFilter(l.ctx, relationDB.DeviceFilter{ProductID: in.ProductID, DeviceNames: []string{in.DeviceName}, WithProduct: true})
		if err != nil {
			return nil, err
		}
		resp := deviceMsg.NewRespCommonMsg(l.ctx, deviceMsg.GetStatus, devices.GenMsgToken(l.ctx, l.svcCtx.NodeID))
		resp.Data = map[string]any{"logLevel": di.LogLevel}

		msg := deviceMsg.PublishMsg{
			Handle:     devices.Log,
			Type:       msgSdkLog.TypeUpdate,
			Payload:    resp.AddStatus(errors.OK, false).Bytes(),
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
	err = l.svcCtx.FastEvent.Publish(l.ctx, topics.DmDeviceInfoUpdate, &devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName})
	if err != nil {
		l.Error(err)
	}
	return &dm.Empty{}, nil
}
