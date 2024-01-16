package vidmgrgbsipmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	db "github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrGbsipChannelDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrGbsipChannelDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrGbsipChannelDeleteLogic {
	return &VidmgrGbsipChannelDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除通道
func (l *VidmgrGbsipChannelDeleteLogic) VidmgrGbsipChannelDelete(in *vid.VidmgrGbsipChannelDelete) (*vid.Response, error) {
	// todo: add your logic here and delete this line
	channelRepo := db.NewVidmgrChannelsRepo(l.ctx)
	filter := db.VidmgrChannelsFilter{
		ChannelIDs: []string{in.ChannelID},
	}
	err := channelRepo.DeleteByFilter(l.ctx, filter)
	if err != nil {
		l.Errorf("%s.Delete err=%v", utils.FuncName(), utils.Fmt(err))
		return nil, err
	}
	return &vid.Response{}, nil
}
