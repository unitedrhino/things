package dc

import (
	"context"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/dcsvr/dc"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ManageGroupMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewManageGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) ManageGroupMemberLogic {
	return ManageGroupMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ManageGroupMemberLogic) ManageGroupMember(req types.ManageGroupMemberReq) (*types.GroupMember, error) {
	l.Infof("ManageGroupMember|req=%+v", req)
	dcReq := &dc.ManageGroupMemberReq{
		Opt: req.Opt,
		Info: &dc.GroupMember{
			GroupID     :req.Info.GroupID,              //组id
			MemberID    :req.Info.MemberID,             //成员id
			MemberType  :req.Info.MemberType,           //成员类型:1:设备 2:用户
		},
	}
	resp, err := l.svcCtx.DcRpc.ManageGroupMember(l.ctx, dcReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.ManageGroupMember|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return RPCToApiFmt(resp).(*types.GroupMember), nil
}
