package otataskmanagelogic

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskDeviceIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaTaskDeviceIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskDeviceIndexLogic {
	return &OtaTaskDeviceIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 升级批次详情列表
func (l *OtaTaskDeviceIndexLogic) OtaTaskDeviceIndex(in *dm.OtaTaskDeviceIndexReq) (*dm.OtaTaskDeviceIndexResp, error) {
	// todo: add your logic here and delete this line

	return &dm.OtaTaskDeviceIndexResp{}, nil
}
