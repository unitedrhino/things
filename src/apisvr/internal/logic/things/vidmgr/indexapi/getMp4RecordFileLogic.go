package indexapi

import (
	"context"
	"encoding/json"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

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
	data, err := proxyMediaServer(l.ctx, l.svcCtx, GETMP4RECORDFILE, req.VidmgrID)
	dataRecv := new(types.IndexApiMp4RecordFileResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
