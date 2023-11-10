package info

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.VidmgrInfoCreateReq) error {
	// todo: add your logic here and delete this line
	vidReq := &vid.VidmgrInfo{
		VidmgrName:   req.VidmgrName,
		VidmgrID:     req.VidmgrID,
		VidmgrIpV4:   req.VidmgrIpV4,
		VidmgrPort:   req.VidmgrPort,
		VidmgrType:   req.VidmgrType,
		VidmgrSecret: req.VidmgrSecret,
		VidmgrStatus: req.VidmgrStatus,
		Desc:         utils.ToRpcNullString(req.Desc),
		Tags:         logic.ToTagsMap(req.Tags),
	}

	jsonStr, _ := json.Marshal(vidReq)
	fmt.Println("[airgens]——Create:", jsonStr)

	_, err := l.svcCtx.VidmgrM.VidmgrInfoCreate(l.ctx, vidReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ManageVidmgr req=%v err=%v", utils.FuncName(), req, er)
		return er
	}

	return nil
}
