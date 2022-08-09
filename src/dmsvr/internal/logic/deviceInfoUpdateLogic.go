package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/service/deviceSend"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"strings"

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

func (l *DeviceInfoUpdateLogic) ChangeDevice(old *mysql.DeviceInfo, data *dm.DeviceInfo) {
	if data.Tags != nil {
		tags, err := json.Marshal(data.Tags)
		if err == nil {
			old.Tags = sql.NullString{
				String: string(tags),
				Valid:  true,
			}
		}
	}
	if data.LogLevel != def.UNKNOWN {
		old.LogLevel = data.LogLevel
	}
	if data.Version != nil {
		old.Version = data.Version.GetValue()
	}
}

// 更新设备
func (l *DeviceInfoUpdateLogic) DeviceInfoUpdate(in *dm.DeviceInfo) (*dm.Response, error) {
	di, err := l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(l.ctx, in.ProductID, in.DeviceName)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.NotFind.AddDetailf("not find device|productid=%s|deviceName=%s",
				in.ProductID, in.DeviceName)
		}
		return nil, errors.Database.AddDetail(err)
	}
	l.ChangeDevice(di, in)

	err = l.svcCtx.DeviceInfo.Update(l.ctx, di)
	if err != nil {
		l.Errorf("ModifyDevice|DeviceInfo|Update|err=%+v", err)
		return nil, errors.System.AddDetail(err)
	}
	//通知device log_level
	uuid, _ := uuid.GenerateUUID()
	tmpTopic := fmt.Sprintf("%s/down/update/%s/%s", devices.TopicHeadLog, di.ProductID, di.DeviceName)
	topic, payload := deviceSend.GenThingDeviceRespData(deviceSend.GET_STATUS, uuid, strings.Split(tmpTopic, "/"), errors.OK, map[string]interface{}{"log_level": di.LogLevel})
	er := l.svcCtx.InnerLink.PublishToDev(l.ctx, topic, payload)
	if er != nil {
		l.Errorf("DeviceResp|SDK Log PublishToDev failure err:%v", er)
	}
	return &dm.Response{}, nil
}
