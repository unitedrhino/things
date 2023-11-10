package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

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

func (l *UpdateLogic) Update(req *types.AreaInfo) error {
	dmReq := &sys.AreaInfo{
		AreaID:       req.AreaID,
		ParentAreaID: req.ParentAreaID,
		ProjectID:    req.ProjectID,
		AreaName:     req.AreaName,
		Position:     logic.ToSysPointRpc(req.Position),
		Desc:         utils.ToRpcNullString(req.Desc),
	}
	_, err := l.svcCtx.AreaM.AreaInfoUpdate(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.AreaManage req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
