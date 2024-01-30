package area

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiDeleteLogic {
	return &MultiDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiDeleteLogic) MultiDelete(req *types.DataAreaMultiDeleteReq) error {
	dto := &sys.DataAreaMultiDeleteReq{
		ProjectID:  req.ProjectID,
		TargetID:   req.TargetID,
		TargetType: req.TargetType,
		AreaIDs:    req.AreaIDs,
	}
	_, err := l.svcCtx.DataM.DataAreaMultiDelete(l.ctx, dto)
	if err != nil {
		l.Errorf("%s.rpc.UserDataAuthManage req=%v err=%v", utils.FuncName(), req, err)
		return err
	}
	return nil
}
