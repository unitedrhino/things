package index

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRtpInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRtpInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRtpInfoLogic {
	return &GetRtpInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRtpInfoLogic) GetRtpInfo(req *types.IndexApiReq) (resp *types.IndexApiRtpInfoResp, err error) {
	// todo: add your logic here and delete this line
	bytetmp := make([]byte, 0)
	data, err := proxySetMediaServer(l.ctx, GETRTPINFO, req.VidmgrID, bytetmp)
	dataRecv := new(types.IndexApiRtpInfoResp)
	json.Unmarshal(data, dataRecv)
	return dataRecv, err
}
