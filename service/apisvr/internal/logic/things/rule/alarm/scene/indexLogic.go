package scene

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

func (l *IndexLogic) Index(req *types.AlarmSceneIndexReq) (resp *types.AlarmSceneMultiSaveReq, err error) {
	ret, err := l.svcCtx.Rule.AlarmSceneIndex(l.ctx, &ud.AlarmSceneIndexReq{
		AlarmID: req.AlarmID,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s req=%v err=%v", utils.FuncName(), req, er)
		return nil, er
	}
	return utils.Copy[types.AlarmSceneMultiSaveReq](ret), nil
}
