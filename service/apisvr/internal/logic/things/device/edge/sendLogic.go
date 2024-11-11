package edge

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/domain/deviceAuth"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dgsvr/pb/dg"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"net/http"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 设备使用http协议用云端交互,需要在http头中带上mqtt的账号密码(basic auth)
func NewSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendLogic {
	return &SendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendLogic) Send(r *http.Request, body []byte, req *types.DeviceInteractEdgeSendReq) (resp any, err error) {
	l.Info(req)
	u, p, ok := r.BasicAuth()
	if !ok {
		return nil, errors.Parameter.WithMsg("账号密码不正确")
	}
	lg, err := deviceAuth.GetLoginDevice(u)
	if err != nil {
		return nil, err
	}
	_, er := l.svcCtx.DeviceA.LoginAuth(l.ctx, &dg.LoginAuthReq{Username: u, //用户名
		Password: p,           //密码
		ClientID: lg.ClientID, //clientID
	})
	if er != nil {
		return nil, er
	}
	if req.Handle[0] == '$' {
		req.Handle = req.Handle[1:]
	}
	_, err = l.svcCtx.DeviceA.AccessAuth(l.ctx, &dg.AccessAuthReq{
		Username: u,
		Topic:    fmt.Sprintf("$%s/up/%s/%s/%s", req.Handle, req.Type, req.ProductID, req.DeviceName),
		ClientID: lg.ClientID,
		Access:   devices.Pub,
	})
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.DeviceInteract.EdgeSend(l.ctx, &dm.EdgeSendReq{
		Handle:     req.Handle,
		Type:       req.Type,
		Payload:    body,
		ProductID:  req.ProductID,
		DeviceName: req.DeviceName,
	})
	if err != nil {
		return nil, err
	}
	var retMap = map[string]any{}
	if len(ret.Payload) > 0 {
		err = json.Unmarshal(ret.Payload, &retMap)
		if err != nil {
			return nil, err
		}
	}
	return retMap, nil
}
