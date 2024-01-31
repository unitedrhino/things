package sipmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/vidsip/internal/media"
	db "github.com/i-Things/things/service/vidsip/internal/repo/relationDB"
	"github.com/i-Things/things/service/vidsip/internal/svc"
	"github.com/i-Things/things/service/vidsip/pb/sip"
	"github.com/zeromicro/go-zero/core/logx"
)

type SipChannelCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSipChannelCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SipChannelCreateLogic {
	return &SipChannelCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 新建通道
func (l *SipChannelCreateLogic) SipChannelCreate(in *sip.SipChnCreateReq) (*sip.Response, error) {
	// todo: add your logic here and delete this line
	// todo: add your logic here and delete this line
	//先判断device ID是否存在
	deviceRepo := db.NewSipDevicesRepo(l.ctx)
	filter := db.SipDevicesFilter{
		DeviceIDs: []string{in.DeviceID},
	}
	device, err := deviceRepo.FindOneByFilter(l.ctx, filter)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s req=%v err=%v", utils.FuncName(), device.DeviceID, er)
		return nil, errors.MediaSipChnCreateError.AddDetail("DeviceID:", device.DeviceID, " 设备不存在")
	}
	channelRepo := db.NewSipChannelsRepo(l.ctx)
	channel := &db.SipChannels{
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
		return nil, errors.MediaSipChnCreateError.AddDetail(er)
	}
	return &sip.Response{}, nil
}
