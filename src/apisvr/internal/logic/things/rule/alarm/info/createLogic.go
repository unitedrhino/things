package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.AlarmInfoCreateReq) (*types.CommonResp, error) {
	rst, err := l.svcCtx.Alarm.AlarmInfoCreate(l.ctx, &rule.AlarmInfo{
		Name:  req.Name,
		State: req.State,
		Desc:  req.Desc,
		Level: req.Level,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.AlarmInfoCreate req=%v err=%v", utils.FuncName(), req, er)
		return nil, er
	}

	return &types.CommonResp{ID: rst.Id}, nil
}
