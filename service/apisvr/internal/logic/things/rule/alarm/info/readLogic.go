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

func (l *ReadLogic) Read(req *types.WithID) (resp *types.AlarmInfo, err error) {
	rpcResp, err := l.svcCtx.Rule.AlarmInfoRead(l.ctx, &ud.WithID{Id: req.ID})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s rpc.AlarmInfoRead req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return AlarmInfoToApi(rpcResp), nil

}
