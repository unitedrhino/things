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

type VidmgrGbsipChannelStopLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrGbsipChannelStopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrGbsipChannelStopLogic {
	return &VidmgrGbsipChannelStopLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 暂停通道
func (l *VidmgrGbsipChannelStopLogic) VidmgrGbsipChannelStop(in *vid.VidmgrGbsipChannelStop) (*vid.Response, error) {
	// todo: add your logic here and delete this line
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
	//err = media.SipStop(po)
	//if err != nil {
	//	return nil, err
	//}

	po.IsPlay = false
	if err := channelRepo.Update(l.ctx, po); err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s channel_id=%v  err=%v", utils.FuncName(), po.ChannelID, er)
		return nil, er
	}
	//update data
	return &vid.Response{}, nil
}
