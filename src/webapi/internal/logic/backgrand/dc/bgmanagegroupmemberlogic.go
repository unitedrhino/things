package dc

import (
	"context"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/webapi/internal/dto"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

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
