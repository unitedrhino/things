package gbsip

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

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
	vidReq := &vid.VidmgrGbsipChannelCreate{
		ChannelID:  req.ChannelID,
		DeviceID:   req.DeviceID,
		Memo:       req.Memo,
		StreamType: req.StreamType,
		Url:        req.URL,
	}
	jsonStr, _ := json.Marshal(req)
	fmt.Println("airgens Createchn :", string(jsonStr))
	_, err := l.svcCtx.VidmgrG.VidmgrGbsipChannelCreate(l.ctx, vidReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ManageVidmgr req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
