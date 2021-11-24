package dc

import (
	"context"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type BgGetGroupMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBgGetGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) BgGetGroupMemberLogic {
	return BgGetGroupMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BgGetGroupMemberLogic) BgGetGroupMember(req types.GetGroupMemberReq) (*types.GetGroupMemberResp, error) {
	// todo: add your logic here and delete this line

	return &types.GetGroupMemberResp{}, nil
}
