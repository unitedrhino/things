package deviceauthlogic

import (
	"context"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/domain/deviceAuth"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccessAuthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	ld *deviceAuth.LoginDevice
	ti *devices.TopicInfo
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
)

//SubSetAuth 网关代理子设备topic校验
func (l *AccessAuthLogic) SubSetAuth(in *dm.AccessAuthReq,
	ld *deviceAuth.LoginDevice, ti *devices.TopicInfo) (err error) {
	_, err = l.svcCtx.Gateway.FindOneByGatewayProductIDGatewayDeviceNameProductIDDeviceName(
		l.ctx, ld.ProductID, ld.DeviceName, ti.ProductID, ti.DeviceName)
	if err != nil {
		if err != mysql.ErrNotFound {
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

//DeviceSelfAuth 设备自己的topic认证
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
	topicInfo, err := devices.GetTopicInfo(in.Topic)
	if err != nil {
		return errors.Permissions
	}
	if ld.ProductID == topicInfo.ProductID && ld.DeviceName == topicInfo.DeviceName {
		return l.DeviceSelfAuth(in)
	}
	return l.SubSetAuth(in, ld, topicInfo)
}

//AccessAuth 设备操作认证
func (l *AccessAuthLogic) AccessAuth(in *dm.AccessAuthReq) (*dm.Response, error) {
	l.Infof("%s req=%+v", utils.FuncName(), utils.Fmt(in))
	err := l.DeviceSelfAuth(in)
	if err != nil {
		l.Infof("%s auth failure=%v", utils.FuncName(), err)
	}
	return &dm.Response{}, err
}
