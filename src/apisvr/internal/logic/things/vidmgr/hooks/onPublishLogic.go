package hooks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OnPublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnPublishLogic {
	return &OnPublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnPublishLogic) OnPublish(req *types.HooksApiPublishReq) (resp *types.HooksApiPublishResp, err error) {
	// todo: add your logic here and delete this line
	reqStr, _ := json.Marshal(*req)

	fmt.Println("---------OnPublish--------------:", string(reqStr))
	return &types.HooksApiPublishResp{
		Code:           0,
		AddMuteAudio:   true,
		ContinuePushMs: 10000,
		EnableAudio:    true,
		EnableFmp4:     true,
		EnableHls:      true,
		EnableHlsFmp4:  true,
		EnableRtmp:     true,
		EnableRtsp:     true,
		EnableTs:       true,
		HlsSavePath:    "/hls_save/path/",
		ModifyStamp:    false,
		Mp4AsPlayer:    false,
		Mp4MaxSecond:   3600,
		Mp4SavePath:    "/mp4_save_path/",
		AutoClose:      false,
		StreamReplace:  "",
	}, nil
}
