package auth

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreaMultiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAreaMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaMultiUpdateLogic {
	return &AreaMultiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AreaMultiUpdateLogic) AreaMultiUpdate(req *types.UserAuthAreaMultiUpdateReq) error {
	dto := &sys.UserAreaMultiUpdateReq{
		UserID:    req.UserID,
		ProjectID: req.ProjectID,
		Areas:     ToAreaPbs(req.Areas),
	}
	_, err := l.svcCtx.UserRpc.UserAreaMultiUpdate(l.ctx, dto)
	if err != nil {
		l.Errorf("%s.rpc.UserDataAuthManage req=%v err=%v", utils.FuncName(), req, err)
		return err
	}
	return nil
}
