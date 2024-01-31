package info

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.WithID) (resp *types.SceneInfo, err error) {
	ruleResp, err := l.svcCtx.Rule.SceneInfoRead(l.ctx, &ud.WithID{Id: req.ID})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s rpc.SceneInfoRead req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return ToSceneTypes(ruleResp), nil
}
