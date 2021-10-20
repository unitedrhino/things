package dc

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &types.GroupMember{}, nil
}
