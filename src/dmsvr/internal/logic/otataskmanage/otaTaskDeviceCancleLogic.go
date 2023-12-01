package otataskmanagelogic

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskDeviceCancleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaTaskDeviceCancleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskDeviceCancleLogic {
	return &OtaTaskDeviceCancleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消单个设备的升级
func (l *OtaTaskDeviceCancleLogic) OtaTaskDeviceCancle(in *dm.OtaTaskDeviceCancleReq) (*dm.OtaCommonResp, error) {
	// todo: add your logic here and delete this line

	return &dm.OtaCommonResp{}, nil
}
