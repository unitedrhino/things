package info

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"gitee.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.DeviceInfo) error {
	in := ToRpcDeviceInfo(req)
	in.Status = 0
	in.Rssi = nil
	in.ProjectID = 0
	in.LogLevel = 0
	in.Version = nil
	in.SoftInfo = ""
	in.HardInfo = ""
	in.IsEnable = 0
	in.IsOnline = 0
	in.MobileOperator = 0
	in.Phone = nil
	in.Iccid = nil
	_, err := l.svcCtx.DeviceM.DeviceInfoUpdate(l.ctx, in)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ManageDevice req=%v err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
