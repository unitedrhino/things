package log

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
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
