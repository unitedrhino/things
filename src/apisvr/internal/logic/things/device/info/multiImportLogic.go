package info

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiImportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiImportLogic {
	return &MultiImportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiImportLogic) MultiImport(req *types.DeviceMultiImportReq) (resp *types.DeviceMultiImportResp, err error) {
	// todo: add your logic here and delete this line

	return
}
