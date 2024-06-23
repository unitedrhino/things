package log

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
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

func (l *IndexLogic) Index(req *types.SceneLogIndexReq) (resp *types.SceneLogIndexResp, err error) {
	ruleResp, err := l.svcCtx.Rule.SceneLogIndex(l.ctx, utils.Copy[ud.SceneLogIndexReq](req))
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.SceneLogIndexReq req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.SceneLogIndexResp{
		Total: ruleResp.Total,
		List:  utils.CopySlice[types.SceneLog](ruleResp.List),
	}, nil
}
