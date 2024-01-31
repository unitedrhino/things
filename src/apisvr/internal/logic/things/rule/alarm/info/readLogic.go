package info

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

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

func (l *ReadLogic) Read(req *types.AlarmInfoReadReq) (resp *types.AlarmInfo, err error) {
	rpcResp, err := l.svcCtx.Alarm.AlarmInfoRead(l.ctx, &rule.WithID{Id: req.ID})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s rpc.AlarmInfoRead req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return AlarmInfoToApi(rpcResp), nil

}
