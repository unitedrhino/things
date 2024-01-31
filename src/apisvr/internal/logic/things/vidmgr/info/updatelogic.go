package info

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.VidmgrInfoUpdateReq) error {
	// todo: add your logic here and delete this line
	vidReq := &vid.VidmgrInfo{
		VidmgrName:   req.VidmgrName,
		VidmgrID:     req.VidmgrID,
		VidmgrIpV4:   req.VidmgrIpV4,
		VidmgrPort:   req.VidmgrPort,
		VidmgrType:   req.VidmgrType,
		VidmgrSecret: req.VidmgrSecret,
		Desc:         utils.ToRpcNullString(req.Desc),
		Tags:         logic.ToTagsMap(req.Tags),
	}
	_, err := l.svcCtx.VidmgrM.VidmgrInfoUpdate(l.ctx, vidReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ManageVidmgr req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
