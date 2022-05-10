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

type BgManageGroupMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBgManageGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) BgManageGroupMemberLogic {
	return BgManageGroupMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BgManageGroupMemberLogic) BgManageGroupMember(req types.ManageGroupMemberReq) (*types.GroupMember, error) {
	l.Infof("ManageGroupMember|req=%+v", req)
	dcReq, err := assemble.ManageGroupMemberReqToRpc(&req)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.DcRpc.ManageGroupMember(l.ctx, dcReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.ManageGroupMember|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return assemble.GroupMemberToApi(resp), nil
}
