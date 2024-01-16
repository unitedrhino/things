package gbsip

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadchnLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadchnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadchnLogic {
	return &ReadchnLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadchnLogic) Readchn(req *types.VidmgrSipReadChnReq) (resp *types.VidmgrSipReadChnResp, err error) {
	// todo: add your logic here and delete this line
	vidResp, err := l.svcCtx.VidmgrG.VidmgrGbsipChannelRead(l.ctx, &vid.VidmgrGbsipChannelRead{
		ChannelID: req.ChannelID,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s rpc.ManageVidmgr req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	apiResp := &types.VidmgrSipReadChnResp{
		CommonSipChannel: *VidmgrGbsipChanneloApi(vidResp),
	}
	return apiResp, nil
}
