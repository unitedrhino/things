package logic

import (
	"context"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetDeviceInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeviceInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetDeviceInfoLogic {
	return GetDeviceInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDeviceInfoLogic) GetDeviceInfo(req types.GetDeviceInfoReq) (*types.GetDeviceInfoResp, error) {
	// todo: add your logic here and delete this line

	return &types.GetDeviceInfoResp{}, nil
}
