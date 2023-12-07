package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type OnShellLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnShellLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnShellLoginLogic {
	return &OnShellLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnShellLoginLogic) OnShellLogin(req *types.HooksApiShellLoginReq) (resp *types.HooksApiResp, err error) {
	// todo: add your logic here and delete this line
	reqStr, _ := json.Marshal(*req)

	fmt.Println("---------OnShellLogin--------------:", string(reqStr))
	return &types.HooksApiResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
