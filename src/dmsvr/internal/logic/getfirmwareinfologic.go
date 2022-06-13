package logic

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFirmwareInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFirmwareInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFirmwareInfoLogic {
	return &GetFirmwareInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取产品固件信息
func (l *GetFirmwareInfoLogic) GetFirmwareInfo(in *dm.GetFirmwareInfoReq) (*dm.GetFirmwareInfoResp, error) {
	// todo: add your logic here and delete this line

	return &dm.GetFirmwareInfoResp{}, nil
}
