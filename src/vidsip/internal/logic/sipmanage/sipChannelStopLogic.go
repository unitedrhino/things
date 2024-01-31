package sipmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	db "github.com/i-Things/things/src/vidsip/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsip/internal/svc"
	"github.com/i-Things/things/src/vidsip/pb/sip"

	"github.com/zeromicro/go-zero/core/logx"
)

type SipChannelStopLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSipChannelStopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SipChannelStopLogic {
	return &SipChannelStopLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 暂停通道
func (l *SipChannelStopLogic) SipChannelStop(in *sip.SipChnStopReq) (*sip.Response, error) {
	// todo: add your logic here and delete this line
	// todo: add your logic here and delete this line
	channelRepo := db.NewSipChannelsRepo(l.ctx)
	filter := db.SipChannelsFilter{ChannelIDs: []string{in.ChannelID}}

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
	return &sip.Response{}, nil
}
