package vidmgrgbsipmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/vidsvr/internal/common"
	db "github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrGbsipChannelReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrGbsipChannelReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrGbsipChannelReadLogic {
	return &VidmgrGbsipChannelReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取通道详情
func (l *VidmgrGbsipChannelReadLogic) VidmgrGbsipChannelRead(in *vid.VidmgrGbsipChannelRead) (*vid.VidmgrGbsipChannel, error) {
	// todo: add your logic here and delete this line
	channelRepo := db.NewVidmgrChannelsRepo(l.ctx)
	filter := db.VidmgrChannelsFilter{
		ChannelIDs: []string{in.ChannelID},
	}
	po, err := channelRepo.FindOneByFilter(l.ctx, filter)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddDetail("Channel not find ID:" + string(in.ChannelID))
		}
		return nil, err
	}

	return common.ToVidmgrGbsipChannelRpc(po), nil
}
