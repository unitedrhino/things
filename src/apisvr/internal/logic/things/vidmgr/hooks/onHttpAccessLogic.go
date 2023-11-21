package hooks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OnHttpAccessLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnHttpAccessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnHttpAccessLogic {
	return &OnHttpAccessLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnHttpAccessLogic) OnHttpAccess(req *types.HooksApiHttpAccessReq) (resp *types.HooksApiResp, err error) {
	// todo: add your logic here and delete this line
	reqStr, _ := json.Marshal(*req)

	fmt.Println("---------OnHttpAccess--------------:", string(reqStr))
	return &types.HooksApiResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
