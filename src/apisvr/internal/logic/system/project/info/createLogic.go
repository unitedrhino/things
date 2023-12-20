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

func (l *CreateLogic) Create(req *types.ProjectInfo) (*types.ProjectWithID, error) {
	rpcReq := &sys.ProjectInfo{
		ProjectName: req.ProjectName,
		CompanyName: utils.ToRpcNullString(req.CompanyName),
		UserID:      req.UserID,
		Region:      utils.ToRpcNullString(req.Region),
		Address:     utils.ToRpcNullString(req.Address),
		Desc:        utils.ToRpcNullString(req.Desc),
	}
	resp, err := l.svcCtx.ProjectM.ProjectInfoCreate(l.ctx, rpcReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ProjectManage req=%v err=%v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.ProjectWithID{ProjectID: resp.ProjectID}, nil
}
