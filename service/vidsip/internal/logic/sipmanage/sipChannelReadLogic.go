package sipmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/vidsip/internal/logic/common"
	db "github.com/i-Things/things/service/vidsip/internal/repo/relationDB"

	"github.com/i-Things/things/service/vidsip/internal/svc"
	"github.com/i-Things/things/service/vidsip/pb/sip"

	"github.com/zeromicro/go-zero/core/logx"
)

type SipChannelReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSipChannelReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SipChannelReadLogic {
	return &SipChannelReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取通道详情
func (l *SipChannelReadLogic) SipChannelRead(in *sip.SipChnReadReq) (*sip.SipChannel, error) {
	// todo: add your logic here and delete this line
	channelRepo := db.NewSipChannelsRepo(l.ctx)
	filter := db.SipChannelsFilter{
		ChannelIDs: []string{in.ChannelID},
	}
	po, err := channelRepo.FindOneByFilter(l.ctx, filter)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddDetail("Channel not find ID:" + string(in.ChannelID))
		}
		return nil, err
	}

	return common.ToSipChannelRpc(po), nil
	//return &sip.SipChannel{}, nil
}
