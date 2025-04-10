package user

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/userSubscribe"
	"github.com/spf13/cast"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubscribeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 插槽用户订阅
func NewSubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscribeLogic {
	return &SubscribeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubscribeLogic) Subscribe(req *types.SlotUserSubscribe) (resp *types.SlotUserSubscribeResp, err error) {
	l.Infof("userSubscribeSlot:%v", utils.Fmt(req))
	resp = &types.SlotUserSubscribeResp{List: []map[string]any{req.Params}}
	switch req.Code {
	case userSubscribe.DeviceConn:
	case userSubscribe.DevicePropertyReport, userSubscribe.DevicePropertyReport2:
	case userSubscribe.DeviceActionReport:
	case userSubscribe.DeviceOtaReport:
	default:
		return resp, errors.NotRealize
	}
	var productID = cast.ToString(req.Params["productID"])
	var deviceName = cast.ToString(req.Params["deviceName"])
	var shareDevice = cast.ToBool(req.Params["shareDevice"])
	var projectID = cast.ToInt64(req.Params["projectID"])
	var areaID = cast.ToInt64(req.Params["areaID"])
	if productID != "" && deviceName != "" { //设备有权限查就能订阅
		_, err := l.svcCtx.DeviceM.DeviceInfoRead(l.ctx, &dm.DeviceInfoReadReq{
			ProductID:  productID,
			DeviceName: deviceName,
		})
		return resp, err
	}
	if shareDevice {
		devs, err := l.svcCtx.DeviceM.DeviceInfoIndex(l.ctx, &dm.DeviceInfoIndexReq{WithShared: 2})
		if err != nil {
			return nil, err
		}
		if len(devs.List) == 0 {
			return resp, errors.NotFind.AddMsg("没有分享的设备")
		}
		resp.List = make([]map[string]interface{}, 0, len(devs.List))
		for _, v := range devs.List {
			resp.List = append(resp.List, map[string]interface{}{"productID": v.ProductID, "deviceName": v.DeviceName})
		}
		return resp, nil
	}
	var uc = ctxs.GetUserCtx(l.ctx)
	if uc == nil {
		return resp, errors.Permissions.AddMsg("只有用户才能订阅")
	}
	if projectID == 0 {
		return resp, errors.Permissions
	}
	if uc.IsAdmin {
		return resp, nil
	}

	pa := uc.ProjectAuth[projectID]
	if pa == nil {
		return resp, errors.Permissions
	}
	if pa.AuthType == def.AuthAdmin || pa.AuthType == def.AuthReadWrite {
		return resp, nil
	}
	if areaID == 0 {
		return resp, errors.Permissions
	}
	aa := pa.Area[areaID]
	if aa == 0 {
		return resp, errors.Permissions
	}
	return resp, nil
}
