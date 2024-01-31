package sipmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/utils"
	db "github.com/i-Things/things/service/vidsip/internal/repo/relationDB"
	"github.com/i-Things/things/service/vidsip/internal/svc"
	"github.com/i-Things/things/service/vidsip/pb/sip"

	"github.com/zeromicro/go-zero/core/logx"
)

type SipChannelDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSipChannelDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SipChannelDeleteLogic {
	return &SipChannelDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除通道
func (l *SipChannelDeleteLogic) SipChannelDelete(in *sip.SipChnDeleteReq) (*sip.Response, error) {
	// todo: add your logic here and delete this line
	channelRepo := db.NewSipChannelsRepo(l.ctx)
	filter := db.SipChannelsFilter{
		ChannelIDs: []string{in.ChannelID},
	}
	err := channelRepo.DeleteByFilter(l.ctx, filter)
	if err != nil {
		l.Errorf("%s.Delete err=%v", utils.FuncName(), utils.Fmt(err))
		return nil, err
	}
	return &sip.Response{}, nil
}
