package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StartRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStartRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartRecordLogic {
	return &StartRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartRecordLogic) StartRecord(req *types.IndexApiReq) (resp *types.IndexApiStartRecordResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, STARTRECORD, req.VidmgrID)
	dataRecv := new(types.IndexApiStartRecordResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
