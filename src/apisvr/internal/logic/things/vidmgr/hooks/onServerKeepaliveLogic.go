package hooks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OnServerKeepaliveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnServerKeepaliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnServerKeepaliveLogic {
	return &OnServerKeepaliveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnServerKeepaliveLogic) OnServerKeepalive(req *types.HooksApiServerKeepaliveReq) (resp *types.HooksApiResp, err error) {
	// todo: add your logic here and delete this line

	//fmt.Println("________keepalive________:", req.Data, "MediaServerId:", req.MediaServerId)

	reqStr, _ := json.Marshal(*req)
	fmt.Println("---------OnServerKeepalive--------------:", string(reqStr))
	//根据MediaServerId 值 判断流媒体服务器是否在线
	//获取当前数据库中服务器的状态值 ,如果不在线侧更新
	//设置一个超时时间,如果超过这个时间,未收到live包,则更新为下线状态

	return &types.HooksApiResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
