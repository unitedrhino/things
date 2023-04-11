package custom

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

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

func (l *UpdateLogic) Update(req *types.ProductCustom) error {
	dmReq := &dm.ProductCustom{
		ProductID:       req.ProductID,
		TransformScript: utils.ToRpcNullString(req.TransformScript),
		ScriptLang:      req.ScriptLang,
		CustomTopic:     req.CustomTopic,
	}
	_, err := l.svcCtx.ProductM.ProductCustomUpdate(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ProductCustomUpdate req=%v err=%+v", utils.FuncName(), utils.Fmt(req), er)
		return er
	}
	return nil
}
