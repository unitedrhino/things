package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type OnStreamNoneReaderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnStreamNoneReaderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnStreamNoneReaderLogic {
	return &OnStreamNoneReaderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnStreamNoneReaderLogic) OnStreamNoneReader(req *types.HooksApiStreamNoneReaderReq) (resp *types.HooksApiStreamNoneReaderResp, err error) {
	// todo: add your logic here and delete this line
	reqStr, _ := json.Marshal(*req)

	fmt.Println("---------OnStreamNoneReader--------------:", string(reqStr))
	return &types.HooksApiStreamNoneReaderResp{}, nil
}
