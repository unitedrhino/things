package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/src/vidsvr/internal/media"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type OnStreamChangedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnStreamChangedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnStreamChangedLogic {
	return &OnStreamChangedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// hooks来的所有数据放进通道后进行处理
func (l *OnStreamChangedLogic) OnStreamChanged(req *types.HooksApiStreamChangedRep) (resp *types.HooksApiResp, err error) {
	// todo: add your logic here and delete this line
	//需要先判断该流服务是否有注册过，未注册过忽略消息
	fmt.Println("____________onStreamChanged_________")
	tmp, _ := json.Marshal(req)

	fmt.Println(string(tmp))
	lstInfo := media.LastStreamInfo{
		LoginTime: time.Now(),
		Req:       *req,
	}
	//并行接收的数据按队列处理
	media.GetMediaChan().ChangeStream <- lstInfo
	fmt.Println(lstInfo)

	return &types.HooksApiResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
