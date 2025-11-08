package edge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gitee.com/unitedrhino/core/share/middlewares"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dgsvr/pb/dg"
	"gitee.com/unitedrhino/things/service/dmsvr/dmExport"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceAuth"
	"github.com/spf13/cast"

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
	if req.Handle[0] == '$' {
		req.Handle = req.Handle[1:]
	}
	err = l.deviceAuth(r, req)
	if err != nil {
		err = l.userAuth(r, req)
		if err != nil {
			return nil, err
		}
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

func (l *SendLogic) userAuth(r *http.Request, req *types.DeviceInteractEdgeSendReq) (err error) {
	userCtx, err := middlewares.Auth(l.ctx, l.svcCtx.UserM, nil, r)
	if err != nil {
		return err
	}
	userCtx.Os = ctxs.GetHandle(r, "User-Agent")
	userCtx.AcceptLanguage = ctxs.GetHandle(r, "Accept-Language")
	strProjectID := ctxs.GetHandle(r, ctxs.UserProjectID, ctxs.UserProjectID2)
	projectID := cast.ToInt64(strProjectID)
	if projectID == 0 {
		projectID = def.NotClassified
	}
	if projectID > def.NotClassified && !userCtx.IsAdmin && userCtx.ProjectAuth[projectID] == nil {
		return errors.Permissions.AddMsg("无所选项目的权限").AddDetailf(strProjectID)
	}
	//注入 用户信息 到 ctx
	ctx2 := ctxs.SetUserCtx(r.Context(), userCtx)
	r = r.WithContext(ctx2)
	err = dmExport.AccessPerm(ctx2, l.svcCtx.DeviceCache, l.svcCtx.UserShareCache, def.AuthRead,
		devices.Core{ProductID: req.ProductID, DeviceName: req.DeviceName}, "")
	return err
}

func (l *SendLogic) deviceAuth(r *http.Request, req *types.DeviceInteractEdgeSendReq) (err error) {
	u, p, ok := r.BasicAuth()
	if !ok {
		return errors.Parameter.WithMsg("账号密码不正确")
	}
	lg, err := deviceAuth.GetLoginDevice(u)
	if err != nil {
		return err
	}
	_, er := l.svcCtx.DeviceA.LoginAuth(l.ctx, &dg.LoginAuthReq{Username: u, //用户名
		Password: p,           //密码
		ClientID: lg.ClientID, //clientID
	})
	if er != nil {
		return er
	}

	if req.ProductID == "" || req.DeviceName == "" {
		req.ProductID = lg.ProductID
		req.DeviceName = lg.DeviceName
	}
	_, err = l.svcCtx.DeviceA.AccessAuth(l.ctx, &dg.AccessAuthReq{
		Username: u,
		Topic:    fmt.Sprintf("$%s/up/%s/%s/%s", req.Handle, req.Type, req.ProductID, req.DeviceName),
		ClientID: lg.ClientID,
		Access:   devices.Pub,
	})
	if err != nil {
		return err
	}
	return nil
}
