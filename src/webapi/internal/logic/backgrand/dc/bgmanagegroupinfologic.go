package dc

import (
	"context"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
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
	// todo: add your logic here and delete this line

	return &types.GroupInfo{}, nil
}
