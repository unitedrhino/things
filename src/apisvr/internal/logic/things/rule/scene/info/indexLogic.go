package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.SceneInfoIndexReq) (resp *types.SceneInfoIndexResp, err error) {
	pbReq := &rule.SceneInfoIndexReq{
		Page:        logic.ToRulePageRpc(req.Page),
		Name:        req.Name,
		State:       req.State,
		TriggerType: req.TriggerType,
	}
	ruleResp, err := l.svcCtx.Scene.SceneInfoIndex(l.ctx, pbReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.SceneInfoIndexReq req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	pis := make([]*types.SceneInfo, 0, len(ruleResp.List))
	for _, v := range ruleResp.List {
		pi := ToSceneTypes(v)
		pis = append(pis, pi)
	}
	return &types.SceneInfoIndexResp{
		Total: ruleResp.Total,
		List:  pis,
		Num:   int64(len(pis)),
	}, nil
}
