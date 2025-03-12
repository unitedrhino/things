package info

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/things/service/udsvr/pb/ud"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

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
		Code: req.Code,
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
		PageResp: logic.ToPageResp(req.Page, ret.Total),
		List:     pis,
	}, nil

}
