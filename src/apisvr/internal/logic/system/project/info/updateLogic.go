package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
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

func (l *UpdateLogic) Update(req *types.ProjectInfo) error {
	dmReq := &sys.ProjectInfo{
		ProjectID:   req.ProjectID,
		ProjectName: req.ProjectName,
		CompanyName: utils.ToRpcNullString(req.CompanyName),
		UserID:      req.UserID,
		Region:      utils.ToRpcNullString(req.Region),
		Address:     utils.ToRpcNullString(req.Address),
		Desc:        utils.ToRpcNullString(req.Desc),
	}
	_, err := l.svcCtx.ProjectM.ProjectInfoUpdate(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ProjectManage req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
