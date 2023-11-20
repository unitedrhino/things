package hooks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OnStreamNotFoundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnStreamNotFoundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnStreamNotFoundLogic {
	return &OnStreamNotFoundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnStreamNotFoundLogic) OnStreamNotFound(req *types.HooksApiStreamNotFoundReq) (resp *types.HooksApiResp, err error) {
	// todo: add your logic here and delete this line
	reqStr, _ := json.Marshal(*req)

	fmt.Println("---------OnStreamNotFound--------------:", string(reqStr))

	return &types.HooksApiResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
