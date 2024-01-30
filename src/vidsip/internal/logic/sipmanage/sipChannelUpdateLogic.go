package sipmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsip/internal/logic/common"
	db "github.com/i-Things/things/src/vidsip/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsip/internal/svc"
	"github.com/i-Things/things/src/vidsip/pb/sip"

	"github.com/zeromicro/go-zero/core/logx"
)

type SipChannelUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSipChannelUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SipChannelUpdateLogic {
	return &SipChannelUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新通道
func (l *SipChannelUpdateLogic) SipChannelUpdate(in *sip.SipChnUpdateReq) (*sip.Response, error) {
	// todo: add your logic here and delete this line
	// todo: add your logic here and delete this line
	channelRepo := db.NewSipChannelsRepo(l.ctx)
	filter := db.SipChannelsFilter{}
	filter.ChannelIDs = []string{in.ChannelID}

	po, err := channelRepo.FindOneByFilter(l.ctx, filter)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddDetail("Channel not find ID:" + string(in.ChannelID))
		}
		return nil, err
	}
	common.UpdatSipChannelDB(po, in)
	if err := channelRepo.Update(l.ctx, po); err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s channel_id=%v  err=%v", utils.FuncName(), po.ChannelID, er)
		return nil, er
	}
	return &sip.Response{}, nil
}
