package dc

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic/assemble"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BgManageGroupInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBgManageGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) BgManageGroupInfoLogic {
	return BgManageGroupInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BgManageGroupInfoLogic) BgManageGroupInfo(req types.ManageGroupInfoReq) (*types.GroupInfo, error) {
	l.Infof("ManageGroupInfo|req=%+v", req)
	dcReq, err := assemble.ManageGroupInfoReqToRpc(&req)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.DcRpc.ManageGroupInfo(l.ctx, dcReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.ManageGroupInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return assemble.GrouInfoToApi(resp), nil

	return &types.GroupInfo{}, nil
}
