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

type CreatechnLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatechnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatechnLogic {
	return &CreatechnLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatechnLogic) Createchn(req *types.VidmgrSipCreateChnReq) error {
	// todo: add your logic here and delete this line
	vidReq := &sip.SipChnCreateReq{
		ChannelID:  req.ChannelID,
		DeviceID:   req.DeviceID,
		Memo:       req.Memo,
		StreamType: req.StreamType,
		Url:        req.URL,
	}
	jsonStr, _ := json.Marshal(req)
	fmt.Println("airgens Createchn :", string(jsonStr))
	_, err := l.svcCtx.SipRpc.SipChannelCreate(l.ctx, vidReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ManageVidmgr req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
