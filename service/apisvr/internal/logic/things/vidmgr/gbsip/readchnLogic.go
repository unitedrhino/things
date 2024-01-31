package gbsip

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/vidsip/pb/sip"

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
	vidResp, err := l.svcCtx.SipRpc.SipChannelRead(l.ctx, &sip.SipChnReadReq{
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
