package dc

import (
	"context"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type BgGetGroupInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBgGetGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) BgGetGroupInfoLogic {
	return BgGetGroupInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BgGetGroupInfoLogic) BgGetGroupInfo(req types.GetGroupInfoReq) (*types.GetGroupInfoResp, error) {
	// todo: add your logic here and delete this line

	return &types.GetGroupInfoResp{}, nil
}
