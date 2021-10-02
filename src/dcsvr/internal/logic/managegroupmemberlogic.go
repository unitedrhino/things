package logic

import (
	"context"

	"gitee.com/godLei6/things/src/dcsvr/dc"
	"gitee.com/godLei6/things/src/dcsvr/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type ManageGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewManageGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManageGroupMemberLogic {
	return &ManageGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 管理组成员
func (l *ManageGroupMemberLogic) ManageGroupMember(in *dc.ManageGroupMemberReq) (*dc.GroupMember, error) {


	return &dc.GroupMember{}, nil
}
