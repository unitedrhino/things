package otataskmanagelogic

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskDeviceProcessLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaTaskDeviceProcessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskDeviceProcessLogic {
	return &OtaTaskDeviceProcessLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 升级进度上报
func (l *OtaTaskDeviceProcessLogic) OtaTaskDeviceProcess(in *dm.OtaTaskDeviceProcessReq) (*dm.OtaCommonResp, error) {
	// todo: add your logic here and delete this line

	return &dm.OtaCommonResp{}, nil
}
