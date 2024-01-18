package project

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectMultiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProjectMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectMultiUpdateLogic {
	return &ProjectMultiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectMultiUpdateLogic) ProjectMultiUpdate(req *types.UserProjectMultiUpdateReq) error {
	dto := &sys.UserProjectMultiUpdateReq{
		UserID:   req.UserID,
		Projects: ToProjectPbs(req.Projects),
	}
	_, err := l.svcCtx.UserRpc.UserProjectMultiUpdate(l.ctx, dto)
	if err != nil {
		l.Errorf("%s.rpc.UserProjectMultiUpdate req=%v err=%v", utils.FuncName(), req, err)
		return err
	}
	return nil
}
