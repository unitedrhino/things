package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type OnRtspRealmLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnRtspRealmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnRtspRealmLogic {
	return &OnRtspRealmLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnRtspRealmLogic) OnRtspRealm(req *types.HooksApiRtspRealmReq) (resp *types.HooksApiRtspRealmResp, err error) {
	// todo: add your logic here and delete this line
	reqStr, _ := json.Marshal(*req)

	fmt.Println("---------OnRtspRealm--------------:", string(reqStr))
	return &types.HooksApiRtspRealmResp{}, nil
}
