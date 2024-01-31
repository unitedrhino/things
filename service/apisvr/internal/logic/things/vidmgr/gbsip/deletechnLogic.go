package gbsip

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/vidsip/pb/sip"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletechnLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletechnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletechnLogic {
	return &DeletechnLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletechnLogic) Deletechn(req *types.VidmgrSipDeleteChnReq) error {
	// todo: add your logic here and delete this line
	vidReq := &sip.SipChnDeleteReq{
		ChannelID: req.ChannelID,
	}
	jsonStr, _ := json.Marshal(req)
	fmt.Println("airgens Deletedev:", string(jsonStr))
	_, err := l.svcCtx.SipRpc.SipChannelDelete(l.ctx, vidReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.Deletechn req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
