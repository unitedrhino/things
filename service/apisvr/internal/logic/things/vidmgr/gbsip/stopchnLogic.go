package gbsip

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/vidsip/pb/sip"

	"github.com/zeromicro/go-zero/core/logx"
)

type StopchnLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStopchnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StopchnLogic {
	return &StopchnLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StopchnLogic) Stopchn(req *types.VidmgrSipStopChnReq) error {
	// todo: add your logic here and delete this line
	// todo: add your logic here and delete this line
	vidReq := &sip.SipChnStopReq{
		ChannelID: req.ChannelID,
	}

	jsonStr, _ := json.Marshal(req)
	fmt.Println("airgens Deletedev:", string(jsonStr))
	_, err := l.svcCtx.SipRpc.SipChannelStop(l.ctx, vidReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.Deletechn req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
