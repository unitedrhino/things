package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
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
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, STARTRECORD, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiStartRecordResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
