package dc

import (
	"context"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ManageGroupInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewManageGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) ManageGroupInfoLogic {
	return ManageGroupInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ManageGroupInfoLogic) ManageGroupInfo(req types.ManageGroupInfoReq) (*types.GroupInfo, error) {
	// todo: add your logic here and delete this line

	return &types.GroupInfo{}, nil
}
