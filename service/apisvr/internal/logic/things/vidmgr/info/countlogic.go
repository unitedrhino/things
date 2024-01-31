package info

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/vidsvr/pb/vid"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CountLogic {
	return &CountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CountLogic) Count(req *types.VidmgrCountReq) (resp *types.VidmgrCountResp, err error) {
	// todo: add your logic here and delete this line
	vidReq := &vid.VidmgrInfoCountReq{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	vidResp, err := l.svcCtx.VidmgrM.VidmgrInfoCount(l.ctx, vidReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.VidmgrInfoCount req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.VidmgrCountResp{
		MgrInfoCount: types.VidmgrInfoCount{
			Online:   vidResp.Online,
			Offline:  vidResp.Offline,
			Inactive: vidResp.Inactive,
		},
	}, nil
}
