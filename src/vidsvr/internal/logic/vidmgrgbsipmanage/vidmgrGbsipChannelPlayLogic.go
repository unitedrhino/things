package vidmgrgbsipmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	db "github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrGbsipChannelPlayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrGbsipChannelPlayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrGbsipChannelPlayLogic {
	return &VidmgrGbsipChannelPlayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 播放通道
func (l *VidmgrGbsipChannelPlayLogic) VidmgrGbsipChannelPlay(in *vid.VidmgrGbsipChannelPlay) (*vid.Response, error) {
	// todo: add your logic here and delete this line
	//play handle
	channelRepo := db.NewVidmgrChannelsRepo(l.ctx)
	filter := db.VidmgrChannelsFilter{ChannelIDs: []string{in.ChannelID}}

	po, err := channelRepo.FindOneByFilter(l.ctx, filter)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddDetail("Channel not find ID:" + string(in.ChannelID))
		}
		return nil, err
	}
	//handle
	//err = media.SipPlayOn(po)
	//err = media.ChnPlay(po.ChannelID)
	if err != nil {
		return nil, err
	}

	po.IsPlay = true
	if err := channelRepo.Update(l.ctx, po); err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s  channel_id=%v  err=%v", utils.FuncName(), po.ChannelID, er)
		return nil, er
	}
	//update data
	return &vid.Response{}, nil
}
