package firmwaremanagelogic

import (
	"context"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaFirmwareFileIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaFirmwareFileIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaFirmwareFileIndexLogic {
	return &OtaFirmwareFileIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 附件列表搜索
func (l *OtaFirmwareFileIndexLogic) OtaFirmwareFileIndex(in *dm.OtaFirmwareFileIndexReq) (*dm.OtaFirmwareFileIndexResp, error) {
	// todo: add your logic here and delete this line

	return &dm.OtaFirmwareFileIndexResp{}, nil
}
