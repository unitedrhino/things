package dc

import (
	"context"
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/shared/utils"
	"github.com/go-things/things/src/webapi/internal/dto"

	"github.com/go-things/things/src/webapi/internal/svc"
	"github.com/go-things/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
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
	dcReq, err := dto.ManageGroupMemberReqToRpc(&req)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.DcRpc.ManageGroupMember(l.ctx, dcReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.ManageGroupMember|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return dto.GroupMemberToApi(resp), nil

	return &types.GroupMember{}, nil
}
