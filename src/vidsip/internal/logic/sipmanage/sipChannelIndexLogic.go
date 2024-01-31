package sipmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/def"
	"github.com/i-Things/things/src/vidsip/internal/logic/common"
	db "github.com/i-Things/things/src/vidsip/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsip/internal/svc"
	"github.com/i-Things/things/src/vidsip/pb/sip"

	"github.com/zeromicro/go-zero/core/logx"
)

type SipChannelIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSipChannelIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SipChannelIndexLogic {
	return &SipChannelIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取通道列表
func (l *SipChannelIndexLogic) SipChannelIndex(in *sip.SipChnIndexReq) (*sip.SipChnIndexResp, error) {
	// todo: add your logic here and delete this line
	channelRepo := db.NewSipChannelsRepo(l.ctx)
	filter := db.SipChannelsFilter{
		ChannelIDs: in.ChannelIDs,
	}
	size, err := channelRepo.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	di, err := channelRepo.FindByFilter(l.ctx, filter, common.ToPageInfoWithDefault(in.Page, &def.PageInfo{
		Page: 1, Size: 20,
		Orders: []def.OrderBy{{"created_time", def.OrderDesc}, {"channel_id", def.OrderDesc}},
	}))
	if err != nil {
		return nil, err
	}
	info := make([]*sip.SipChannel, 0, len(di))
	for _, v := range di {
		info = append(info, common.ToSipChannelRpc(v))
	}
	return &sip.SipChnIndexResp{List: info, Total: size}, nil
}
