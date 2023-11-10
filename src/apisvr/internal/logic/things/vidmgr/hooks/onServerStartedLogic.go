package hooks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OnServerStartedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnServerStartedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnServerStartedLogic {
	return &OnServerStartedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnServerStartedLogic) OnServerStarted(req *types.HooksApiServerStartedReq) (resp *types.HooksApiResp, err error) {
	// todo: add your logic here and delete this line
	reqStr, _ := json.Marshal(*req)

	fmt.Println("---------OnServerStarted--------------:", string(reqStr))
	return &types.HooksApiResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
