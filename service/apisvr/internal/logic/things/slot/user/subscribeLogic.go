package user

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/spf13/cast"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubscribeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscribeLogic {
	return &SubscribeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubscribeLogic) Subscribe(req *types.SlotUserSubscribeReq) error {
	l.Infof("userSubscribeSlot:%v", utils.Fmt(req))
	switch req.Code {
	case def.UserSubscribeDeviceConn:
	case def.UserSubscribeDevicePropertyReport:
	case def.UserSubscribeDeviceOtaReport:
	default:
		return errors.NotRealize
	}
	var productID = cast.ToString(req.Params["productID"])
	var deviceName = cast.ToString(req.Params["deviceName"])
	var projectID = cast.ToInt64(req.Params["projectID"])
	var areaID = cast.ToInt64(req.Params["areaID"])
	if productID != "" && deviceName != "" { //设备有权限查就能订阅
		_, err := l.svcCtx.DeviceM.DeviceInfoRead(l.ctx, &dm.DeviceInfoReadReq{
			ProductID:  productID,
			DeviceName: deviceName,
		})
		return err
	}
	var uc = ctxs.GetUserCtx(l.ctx)
	if uc == nil {
		return errors.Permissions.AddMsg("只有用户才能订阅")
	}
	if projectID == 0 {
		return errors.Permissions
	}
	if uc.IsAdmin {
		return nil
	}
	pa := uc.ProjectAuth[projectID]
	if pa == nil {
		return errors.Permissions
	}
	if pa.AuthType == def.AuthAdmin || pa.AuthType == def.AuthReadWrite {
		return nil
	}
	if areaID == 0 {
		return errors.Permissions
	}
	aa := pa.Area[areaID]
	if aa == 0 {
		return errors.Permissions
	}
	return nil
}
