package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type OnRecordMp4Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnRecordMp4Logic(ctx context.Context, svcCtx *svc.ServiceContext) *OnRecordMp4Logic {
	return &OnRecordMp4Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnRecordMp4Logic) OnRecordMp4(req *types.HooksApiRecordMp4Req) (resp *types.HooksApiResp, err error) {
	// todo: add your logic here and delete this line
	reqStr, _ := json.Marshal(*req)

	fmt.Println("---------OnRecordMp4--------------:", string(reqStr))
	return &types.HooksApiResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
