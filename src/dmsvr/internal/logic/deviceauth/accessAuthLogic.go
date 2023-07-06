package deviceauthlogic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/domain/deviceAuth"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsgManage"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccessAuthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	ld   *deviceAuth.LoginDevice
	ti   *devices.TopicInfo
	GdDB *relationDB.GatewayDeviceRepo
}

func NewAccessAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccessAuthLogic {
	return &AccessAuthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		GdDB:   relationDB.NewGatewayDeviceRepo(ctx),
	}
}

var (
	AccessMap = map[string]devices.Direction{
		devices.Pub: devices.Up,
		devices.Sub: devices.Down,
	}
	AccessToActionMap = map[string]string{
		devices.Pub: "publish",
		devices.Sub: "subscribe",
	}
)

// SubSetAuth 网关代理子设备topic校验
func (l *AccessAuthLogic) SubSetAuth(in *dm.AccessAuthReq,
	ld *deviceAuth.LoginDevice, ti *devices.TopicInfo) (err error) {
	_, err = l.GdDB.FindOneByFilter(l.ctx, relationDB.GatewayDeviceFilter{
		Gateway: &devices.Core{
			ProductID:  ld.ProductID,
			DeviceName: ld.DeviceName,
		},
		SubDevice: &devices.Core{
			ProductID:  ti.ProductID,
			DeviceName: ti.DeviceName,
		},
	})
	if err != nil {
		if !errors.Cmp(err, errors.NotFind) {
			return errors.Database.AddDetail(err)
		}
		return errors.Permissions
	}
	access, ok := AccessMap[in.Access]
	if !ok {
		return errors.Permissions
	}
	if access != ti.Direction {
		return errors.Permissions
	}
	return nil
}

// DeviceSelfAuth 设备自己的topic认证
func (l *AccessAuthLogic) DeviceSelfAuth(in *dm.AccessAuthReq) (err error) {
	l.Infof("%s in:%v", utils.FuncName(), utils.Fmt(in))
	access, ok := AccessMap[in.Access]
	if !ok {
		return errors.Permissions
	}
	return deviceAuth.AccessAuth(deviceAuth.AuthInfo{
		Username: in.Username,
		Topic:    in.Topic,
		ClientID: in.ClientID,
		Access:   access,
		Ip:       in.Ip,
	})
}

func (l *AccessAuthLogic) Auth(in *dm.AccessAuthReq) (err error) {
	ld, err := deviceAuth.GetClientIDInfo(in.ClientID)
	if err != nil {
		return err
	}
	err = func() error {
		topicInfo, err := devices.GetTopicInfo(in.Topic)
		if err != nil {
			return errors.Permissions
		}
		if ld.ProductID == topicInfo.ProductID && ld.DeviceName == topicInfo.DeviceName {
			return l.DeviceSelfAuth(in)
		}
		return l.SubSetAuth(in, ld, topicInfo)
	}()
	er := l.svcCtx.HubLogRepo.Insert(l.ctx, &deviceMsgManage.HubLog{
		ProductID:  ld.ProductID,
		Action:     AccessToActionMap[in.Access],
		Timestamp:  time.Now(), // 操作时间
		DeviceName: ld.DeviceName,
		Topic:      in.Topic,
		Content:    fmt.Sprintf("ip:%v", in.Ip),
		TranceID:   utils.TraceIdFromContext(l.ctx),
		ResultType: errors.Fmt(err).GetCode(),
	})
	if er != nil {
		l.Errorf("%v.HubLogRepo.Insert err:%v", utils.FuncName(), er)
	}
	return err
}

// AccessAuth 设备操作认证
func (l *AccessAuthLogic) AccessAuth(in *dm.AccessAuthReq) (*dm.Response, error) {
	l.Infof("%s req=%+v", utils.FuncName(), utils.Fmt(in))
	err := l.Auth(in)
	if err != nil {
		l.Infof("%s auth failure=%v", utils.FuncName(), err)
	}
	return &dm.Response{}, err
}
