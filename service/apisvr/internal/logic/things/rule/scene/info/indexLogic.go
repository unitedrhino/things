package info

import (
	"context"
	"gitee.com/unitedrhino/share/def"
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

func (l *IndexLogic) Index(req *types.SceneInfoIndexReq) (resp *types.SceneInfoIndexResp, err error) {
	ruleResp, err := l.svcCtx.Rule.SceneInfoIndex(l.ctx, utils.Copy[ud.SceneInfoIndexReq](req))
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.SceneInfoIndexReq req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	pis := make([]*types.SceneInfo, 0, len(ruleResp.List))
	for _, v := range ruleResp.List {
		pi := ToSceneTypes(v)
		if req.IsOnlyCore == def.True {
			pi.If = ""
			pi.When = ""
			pi.Then = ""
		}
		pis = append(pis, pi)
	}
	return &types.SceneInfoIndexResp{
		PageResp: logic.ToPageResp(req.Page, ruleResp.Total),
		List:     pis,
		Num:      int64(len(pis)),
	}, nil
}
