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

func (l *IndexLogic) Index(req *types.AlarmInfoIndexReq) (resp *types.AlarmInfoIndexResp, err error) {
	ret, err := l.svcCtx.Alarm.AlarmInfoIndex(l.ctx, &rule.AlarmInfoIndexReq{
		Page:     logic.ToRulePageRpc(req.Page),
		Name:     req.Name,
		SceneID:  req.SceneID,
		AlarmIDs: req.AlarmIDs,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.AlarmInfoIndex req=%v err=%v", utils.FuncName(), req, er)
		return nil, er
	}
	pis := make([]*types.AlarmInfo, 0, len(ret.List))
	for _, v := range ret.List {
		pis = append(pis, AlarmInfoToApi(v))
	}
	return &types.AlarmInfoIndexResp{
		Total: ret.Total,
		List:  pis,
		Num:   int64(len(pis)),
	}, nil

}
