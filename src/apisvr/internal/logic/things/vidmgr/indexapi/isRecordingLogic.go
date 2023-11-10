package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsRecordingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIsRecordingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsRecordingLogic {
	return &IsRecordingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IsRecordingLogic) IsRecording(req *types.IndexApiReq) (resp *types.IndexApiIsRecordingResp, err error) {
	// todo: add your logic here and delete this line
	data, err := proxyMediaServer(l.ctx, l.svcCtx, ISRECORDING, req.VidmgrID)
	dataRecv := new(types.IndexApiIsRecordingResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
