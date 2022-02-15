package logic

import (
	"context"

	"github.com/go-things/things/src/dmsvr/dm"
	"github.com/go-things/things/src/dmsvr/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeviceDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeviceDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceDataLogic {
	return &GetDeviceDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDeviceDataLogic) GetDeviceData(in *dm.GetDeviceDataReq) (*dm.GetDeviceDataResp, error) {
	// todo: add your logic here and delete this line

	return &dm.GetDeviceDataResp{}, nil
}
