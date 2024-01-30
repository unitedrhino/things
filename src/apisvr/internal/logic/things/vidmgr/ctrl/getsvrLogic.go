package ctrl

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetsvrLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetsvrLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetsvrLogic {
	return &GetsvrLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

//1:zlmediakit,2:srs,3:monibuca
func (l *GetsvrLogic) Getsvr(req *types.CtrlApiReq) (resp *types.CtrlApiResp, err error) {
	// todo: add your logic here and delete this line
	resp = &types.CtrlApiResp{
		Code: 0,
	}
	switch req.SrvType {
	case 1:
		byteTmp, errTmp := handleZLMediakitReq(req)
		if errTmp != nil {
			resp.Code = 500
			resp.Data = ""
			err = errTmp
		}
		err = nil
		resp.Data = string(byteTmp)
	case 2:
	case 3:
	default:
	}
	return resp, err
}
