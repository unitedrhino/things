package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMp4RecordFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMp4RecordFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMp4RecordFileLogic {
	return &GetMp4RecordFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMp4RecordFileLogic) GetMp4RecordFile(req *types.IndexApiReq) (resp *types.IndexApiMp4RecordFileResp, err error) {
	// todo: add your logic here and delete this line
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, GETMP4RECORDFILE, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiMp4RecordFileResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
