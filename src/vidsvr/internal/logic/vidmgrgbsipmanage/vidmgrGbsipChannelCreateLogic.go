package vidmgrgbsipmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/media"
	db "github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrGbsipChannelCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrGbsipChannelCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrGbsipChannelCreateLogic {
	return &VidmgrGbsipChannelCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 新建通道
func (l *VidmgrGbsipChannelCreateLogic) VidmgrGbsipChannelCreate(in *vid.VidmgrGbsipChannelCreate) (*vid.Response, error) {
	// todo: add your logic here and delete this line
	//先判断device ID是否存在
	deviceRepo := db.NewVidmgrDevicesRepo(l.ctx)
	filter := db.VidmgrDevicesFilter{
		DeviceIDs: []string{in.DeviceID},
	}
	device, err := deviceRepo.FindOneByFilter(l.ctx, filter)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s req=%v err=%v", utils.FuncName(), device.DeviceID, er)
		return nil, errors.MediaGbsipChnCreateError.AddDetail("DeviceID:", device.DeviceID, " 设备不存在")
	}
	channelRepo := db.NewVidmgrChannelsRepo(l.ctx)
	channel := &db.VidmgrChannels{
		ChannelID: in.ChannelID,
		//ChannelID:  fmt.Sprintf("%s%06d", media.SipInfo.CID, media.SipInfo.CNUM+1),
		DeviceID:   in.DeviceID,
		Memo:       in.Memo,
		StreamType: in.StreamType,
		URL:        in.Url,
	}
	//更新数据
	media.SipInfo.CNUM += 1
	//l.svcCtx.Config.GbsipConf.Cnum = strconv.ParseInt(media.SipInfo.CNUM, 10, 32)

	if err := channelRepo.Insert(l.ctx, channel); err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s req=%v err=%v", utils.FuncName(), channel.ChannelID, er)
		return nil, errors.MediaGbsipChnCreateError.AddDetail(er)
	}
	return &vid.Response{}, nil
}
