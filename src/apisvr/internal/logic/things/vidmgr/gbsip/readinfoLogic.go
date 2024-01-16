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

type ReadinfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadinfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadinfoLogic {
	return &ReadinfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadinfoLogic) Readinfo(req *types.VidmgrSipReadInfoReq) (resp *types.VidmgrSipReadInfoResp, err error) {
	// todo: add your logic here and delete this line
	vidResp, err := l.svcCtx.VidmgrG.VidgmrGbsipInfoRead(l.ctx, &vid.VidmgrGbsipInfoReadReq{
		VidmgrID: req.VidmgrID,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s rpc.ManageVidmgr req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	apiResp := &types.VidmgrSipReadInfoResp{
		CommonSipInfo: *ToVidmgrGbsipInfoApi(vidResp),
	}

	return apiResp, nil
}
