package logic

import (
	"context"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/dmsvr/dm"
	"github.com/golang/protobuf/ptypes/wrappers"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ManageDeviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewManageDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) ManageDeviceLogic {
	return ManageDeviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ManageDeviceLogic) ManageDevice(req types.ManageDeviceReq) (*types.DeviceInfo, error) {
	defer func() {
		if p := recover(); p != nil {
			utils.HandleThrow(p)
		}
	}()
	dmReq := &dm.ManageDeviceReq{
		Opt: req.Opt,
		Info: &dm.DeviceInfo{
			ProductID:  req.Info.ProductID,  //产品id 只读
			DeviceName: req.Info.DeviceName, //设备名称 读写
			LogLevel:   req.Info.LogLevel,   // 日志级别:1)关闭 2)错误 3)告警 4)信息 5)调试  读写
		},
	}
	if req.Info.Version != nil {
		dmReq.Info.Version = &wrappers.StringValue{
			Value: *req.Info.Version,
		}
	}
	resp, err := l.svcCtx.DmRpc.ManageDevice(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.ManageDevice|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return RPCToApiFmt(resp).(*types.DeviceInfo), nil
}
