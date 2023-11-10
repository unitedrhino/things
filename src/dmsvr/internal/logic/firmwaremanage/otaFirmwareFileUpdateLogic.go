package firmwaremanagelogic

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaFirmwareFileUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaFirmwareFileUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaFirmwareFileUpdateLogic {
	return &OtaFirmwareFileUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 附件信息更新
func (l *OtaFirmwareFileUpdateLogic) OtaFirmwareFileUpdate(in *dm.OtaFirmwareFileReq) (*dm.OtaFirmwareFileResp, error) {
	// todo: add your logic here and delete this line

	return &dm.OtaFirmwareFileResp{}, nil
}
