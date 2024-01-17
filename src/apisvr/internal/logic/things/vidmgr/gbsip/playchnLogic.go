package gbsip

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/vidsip/pb/sip"
	"github.com/zeromicro/go-zero/core/logx"
)

type PlaychnLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPlaychnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PlaychnLogic {
	return &PlaychnLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PlaychnLogic) Playchn(req *types.VidmgrSipPlayChnReq) error {
	// todo: add your logic here and delete this line
	fmt.Println("-----------------Playchn----------------------")
	vidReq := &sip.SipChnPlayReq{
		ChannelID: req.ChannelID,
	}
	jsonStr, _ := json.Marshal(req)
	fmt.Println("airgens Deletedev:", string(jsonStr))
	_, err := l.svcCtx.SipRpc.SipChannelPlay(l.ctx, vidReq)
	if err != nil {
		er := errors.Fmt(err)
		fmt.Printf("%s.rpc.Deletechn req=%v err=%v", utils.FuncName(), req, er)
		//return err
	}
	return nil
}
