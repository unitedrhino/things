package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/service/vidsvr/internal/svc"
	"github.com/i-Things/things/service/vidsvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type OnPlayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnPlayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnPlayLogic {
	return &OnPlayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnPlayLogic) OnPlay(req *types.HooksApiPlayReq) (resp *types.HooksApiResp, err error) {
	// todo: add your logic here and delete this line
	reqStr, _ := json.Marshal(*req)

	fmt.Println("---------OnPlay--------------:", string(reqStr))
	return &types.HooksApiResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
