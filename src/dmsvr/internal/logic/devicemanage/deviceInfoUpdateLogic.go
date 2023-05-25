package devicemanagelogic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoUpdateLogic {
	return &DeviceInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceInfoUpdateLogic) SetDevicePoByDto(old *mysql.DmDeviceInfo, data *dm.DeviceInfo) {
	if data.Tags != nil {
		tags, err := json.Marshal(data.Tags)
		if err == nil {
			old.Tags = string(tags)
		}
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
	if data.Version != nil {
		old.Version = data.Version.GetValue()
	}
	if data.HardInfo != "" {
		old.HardInfo = data.HardInfo
	}
	if data.SoftInfo != "" {
		old.SoftInfo = data.SoftInfo
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
		old.Position = fmt.Sprintf("POINT(%f %f)", data.Position.Longitude, data.Position.Latitude)
	}

	if data.DeviceAlias != nil {
		old.DeviceAlias = utils.AnyToNullString(data.DeviceAlias).String
	}
}

// 更新设备
func (l *DeviceInfoUpdateLogic) DeviceInfoUpdate(in *dm.DeviceInfo) (*dm.Response, error) {
	if in.ProductID == "" && in.ProductName != "" { //通过唯一的产品名 查找唯一的产品ID
		if pid, err := l.svcCtx.ProductInfo.GetIDByName(l.ctx, mysql.ProductFilter{ProductName: in.ProductName}, nil); err != nil {
			return nil, err
		} else {
			in.ProductID = pid
		}
	}

	dmDiPo, err := l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(l.ctx, in.ProductID, in.DeviceName)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.NotFind.AddDetailf("not find device productID=%s deviceName=%s",
				in.ProductID, in.DeviceName)
		}
		return nil, errors.Database.AddDetail(err)
	}

	l.SetDevicePoByDto(dmDiPo, in)

	err = l.svcCtx.DeviceInfo.UpdateDeviceInfo(l.ctx, dmDiPo)
	if err != nil {
		l.Errorf("DeviceInfoUpdate.DeviceInfo.Update err=%+v", err)
		return nil, errors.System.AddDetail(err)
	}

	if in.LogLevel != def.Unknown {
		err := l.svcCtx.DataUpdate.DeviceLogLevelUpdate(l.ctx, &events.DeviceUpdateInfo{
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
		})
		if err != nil {
			l.Errorf("DeviceInfoUpdate.DeviceLogLevelUpdate err=%+v", err)
		}
	}
	return &dm.Response{}, nil
}
