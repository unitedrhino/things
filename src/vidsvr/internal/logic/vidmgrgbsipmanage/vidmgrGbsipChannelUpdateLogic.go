package vidmgrgbsipmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/common"
	db "github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrGbsipChannelUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrGbsipChannelUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrGbsipChannelUpdateLogic {
	return &VidmgrGbsipChannelUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新通道
func (l *VidmgrGbsipChannelUpdateLogic) VidmgrGbsipChannelUpdate(in *vid.VidmgrGbsipChannelUpdate) (*vid.Response, error) {
	// todo: add your logic here and delete this line
	channelRepo := db.NewVidmgrChannelsRepo(l.ctx)
	filter := db.VidmgrChannelsFilter{}
	filter.ChannelIDs = []string{in.ChannelID}

	po, err := channelRepo.FindOneByFilter(l.ctx, filter)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddDetail("Channel not find ID:" + string(in.ChannelID))
		}
		return nil, err
	}
	common.UpdatVidmgrChannelDB(po, in)
	if err := channelRepo.Update(l.ctx, po); err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s channel_id=%v  err=%v", utils.FuncName(), po.ChannelID, er)
		return nil, er
	}
	return &vid.Response{}, nil
}
