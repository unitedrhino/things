package stream

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/vidsvr/client/vidmgrinfomanage"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.VidmgrStreamReadReq) (resp *types.VidmgrStreamReadResp, err error) {
	// todo: add your logic here and delete this line
	vidResp, err := l.svcCtx.VidmgrS.VidmgrStreamRead(l.ctx, &vidmgrinfomanage.VidmgrStreamReadReq{
		StreamID: req.StreamID,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s rpc.ManageVidmgr req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	apiResp := &types.VidmgrStreamReadResp{
		VidmgrStream: *VidmgrStreamToApi(vidResp),
		MediaPort:    vidResp.MediaPort,
		MediaIP:      vidResp.MediaIP,
	}

	return apiResp, nil
}
