package info

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/logic"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

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
	ret, err := l.svcCtx.Rule.AlarmInfoIndex(l.ctx, &ud.AlarmInfoIndexReq{
		Page: logic.ToUdPageRpc(req.Page),
		Name: req.Name,
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
	}, nil

}
