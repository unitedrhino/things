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
	// todo: add your logic here and delete this line
	vidReq := &vid.VidmgrGbsipChannelDelete{
		ChannelID: req.ChannelID,
	}
	jsonStr, _ := json.Marshal(req)
	fmt.Println("airgens Deletedev:", string(jsonStr))
	_, err := l.svcCtx.VidmgrG.VidmgrGbsipChannelDelete(l.ctx, vidReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.Deletechn req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
