package job

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FirmwareIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFirmwareIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FirmwareIndexLogic {
	return &FirmwareIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FirmwareIndexLogic) FirmwareIndex(req *types.OtaJobByFirmwareIndexReq) (resp *types.OtaJobInfoIndexResp, err error) {
	// todo: add your logic here and delete this line

	return
}
