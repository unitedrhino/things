package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/service/vidsvr/internal/svc"
	"github.com/i-Things/things/service/vidsvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type OnRtspAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnRtspAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnRtspAuthLogic {
	return &OnRtspAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnRtspAuthLogic) OnRtspAuth(req *types.HooksApiPublishReq) (resp *types.HooksApiRtspAuthResp, err error) {
	// todo: add your logic here and delete this line
	reqStr, _ := json.Marshal(*req)

	fmt.Println("---------OnRtspAuth--------------:", string(reqStr))
	return &types.HooksApiRtspAuthResp{}, nil
}
