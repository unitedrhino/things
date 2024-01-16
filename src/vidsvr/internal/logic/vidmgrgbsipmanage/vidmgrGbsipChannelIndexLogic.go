package vidmgrgbsipmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/vidsvr/internal/common"
	db "github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrGbsipChannelIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrGbsipChannelIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrGbsipChannelIndexLogic {
	return &VidmgrGbsipChannelIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取通道列表
func (l *VidmgrGbsipChannelIndexLogic) VidmgrGbsipChannelIndex(in *vid.VidmgrGbsipChannelIndexReq) (*vid.VidmgrGbsipChannelIndexResp, error) {
	// todo: add your logic here and delete this line
	channelRepo := db.NewVidmgrChannelsRepo(l.ctx)
	filter := db.VidmgrChannelsFilter{
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
	info := make([]*vid.VidmgrGbsipChannel, 0, len(di))
	for _, v := range di {
		info = append(info, common.ToVidmgrGbsipChannelRpc(v))
	}
	return &vid.VidmgrGbsipChannelIndexResp{List: info, Total: size}, nil
}
