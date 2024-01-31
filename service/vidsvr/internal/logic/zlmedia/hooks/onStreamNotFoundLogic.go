package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/service/vidsvr/internal/svc"
	"github.com/i-Things/things/service/vidsvr/internal/types"
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
	fmt.Println("-------------OnStreamNotFound-------------:", string(reqStr))
	//当播放触发notfound时重新拉流
	//查表
	/*
		streamRepo := db.NewVidmgrStreamRepo(l.ctx)
		filter := db.VidmgrStreamFilter{
			Stream: req.Stream,
		}

			stream, err := streamRepo.FindOneByFilter(l.ctx, filter)
			if stream != nil {
				//根据流类型，确定
				if stream.OriginType == clients.RTMP_PUSH || stream.OriginType == clients.RTSP_PUSH ||
					stream.OriginType == clients.RTP_PUSH {
					// 存在推流记录关闭当前，重新发起推流
					media.SipStopPlay(stream.Stream)
				} else { //pull
					// 拉流的，重新拉流
					media.SipPlay(stream)
				}
			}*/
	return &types.HooksApiResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
