package deviceauthlogic

import (
	"context"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/domain/deviceAuth"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

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
)

func (l *AccessAuthLogic) Auth(in *dm.AccessAuthReq) (err error) {
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

// 设备操作认证
func (l *AccessAuthLogic) AccessAuth(in *dm.AccessAuthReq) (*dm.Response, error) {
	l.Infof("%s|req=%+v", utils.FuncName(), utils.Fmt(in))
	err := l.Auth(in)
	if err != nil {
		l.Infof("%s|auth failure=%v", utils.FuncName(), err)
	}
	return &dm.Response{}, err
}
