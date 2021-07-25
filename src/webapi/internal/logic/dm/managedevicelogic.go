package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &types.DeviceInfo{}, nil
}
