package deviceauthlogic

import (
	"context"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/deviceAuth"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/dgsvr/internal/svc"
	"gitee.com/i-Things/things/service/dgsvr/pb/dg"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccessAuthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAccessAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccessAuthLogic {
	return &AccessAuthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
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

// DeviceSelfAuth 设备自己的topic认证
func (l *AccessAuthLogic) DeviceSelfAuth(in *dg.AccessAuthReq) (err error) {
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

// SubSetAuth 网关代理子设备topic校验
func (l *AccessAuthLogic) SubSetAuth(in *dg.AccessAuthReq,
	ld *deviceAuth.LoginDevice, ti *devices.TopicInfo) (err error) {
	ret, err := l.svcCtx.DeviceM.DeviceGatewayIndex(l.ctx, &dm.DeviceGatewayIndexReq{
		Gateway: &dm.DeviceCore{
			ProductID:  ld.ProductID,
			DeviceName: ld.DeviceName,
		},
		SubDevice: &dm.DeviceCore{
			ProductID:  ti.ProductID,
			DeviceName: ti.DeviceName,
		},
	})
	if err != nil {
		return errors.Fmt(err)
	}
	if ret.Total == 0 { //未找到该网关和设备的关系
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

func (l *AccessAuthLogic) Auth(in *dg.AccessAuthReq) (err error) {
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
	return err
}

// 设备操作认证
func (l *AccessAuthLogic) AccessAuth(in *dg.AccessAuthReq) (*dg.Response, error) {
	l.Infof("%s req=%+v", utils.FuncName(), utils.Fmt(in))
	err := l.Auth(in)
	if err != nil {
		l.Infof("%s auth failure=%v", utils.FuncName(), err)
	}
	return &dg.Response{}, err
}
