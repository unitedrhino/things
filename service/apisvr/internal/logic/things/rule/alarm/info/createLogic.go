package info

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/udsvr/pb/ud"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

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

func (l *CreateLogic) Create(req *types.AlarmInfo) (*types.CommonResp, error) {
	rst, err := l.svcCtx.Rule.AlarmInfoCreate(l.ctx, utils.Copy[ud.AlarmInfo](req))
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.AlarmInfoCreate req=%v err=%v", utils.FuncName(), req, er)
		return nil, er
	}

	return &types.CommonResp{ID: rst.Id}, nil
}
