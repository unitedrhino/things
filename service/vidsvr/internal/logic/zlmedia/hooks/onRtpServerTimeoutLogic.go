package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/service/vidsvr/internal/svc"
	"github.com/i-Things/things/service/vidsvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type OnRtpServerTimeoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnRtpServerTimeoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnRtpServerTimeoutLogic {
	return &OnRtpServerTimeoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnRtpServerTimeoutLogic) OnRtpServerTimeout(req *types.HooksApiRtpServerTimeoutReq) (resp *types.HooksApiResp, err error) {
	// todo: add your logic here and delete this line
	reqStr, _ := json.Marshal(*req)

	fmt.Println("---------OnRtpServerTimeout--------------:", string(reqStr))
	return &types.HooksApiResp{}, nil
}
