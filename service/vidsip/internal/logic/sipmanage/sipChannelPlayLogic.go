package sipmanagelogic

import (
	"context"
	"github.com/i-Things/things/service/vidsip/internal/svc"
	"github.com/i-Things/things/service/vidsip/pb/sip"

	"github.com/zeromicro/go-zero/core/logx"
)

type SipChannelPlayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSipChannelPlayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SipChannelPlayLogic {
	return &SipChannelPlayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//play handle
//fmt.Println("-----------------VidmgrGbsipChannelPlay----------------------")
//
//channelRepo := db.NewSipChannelsRepo(l.ctx)
//filter := db.SipChannelsFilter{ChannelIDs: []string{in.ChannelID}}
//
//po, err := channelRepo.FindOneByFilter(l.ctx, filter)
//if err != nil {
//	fmt.Println("-----------------VidmgrGbsipChannelPlay  channelRepo Error----------------------")
//	if errors.Cmp(err, errors.NotFind) {
//		return nil, errors.MediaSipPlayError.AddDetail("Channel not find ID:" + string(in.ChannelID))
//	}
//	return nil, err
//}
////if po.Status != 1 {
////	fmt.Println("-----------------VidmgrGbsipChannelPlay  po.Status != 1----------------------")
////	return nil, errors.MediaSipPlayError.AddDetail("通道未在线")
////}
//deviceRepo := db.NewSipDevicesRepo(l.ctx)
//
//deviceFilter := db.SipDevicesFilter{
//	DeviceID: po.DeviceID,
//}
//pdev, err := deviceRepo.FindOneByFilter(l.ctx, deviceFilter)
//if err != nil {
//	fmt.Println("-----------------VidmgrGbsipChannelPlay  deviceRepo----------------------")
//	if errors.Cmp(err, errors.NotFind) {
//		return nil, errors.MediaSipPlayError.AddDetail("Device not find ID:" + string(pdev.DeviceID))
//	}
//	return nil, err
//}
//
//stream := &media.Stream{
//	Type:       0, //直播类型
//	ChannelID:  po.ChannelID,
//	DeviceID:   po.DeviceID,
//	Stream:     fmt.Sprintf("%d%s%04d", 0, po.DeviceID[15:20], 1),
//	ChnURIStr:  po.URIStr,
//	DevNetType: "udp",
//	DevSource:  pdev.Source,
//	//MediaIP:    utils.InetNtoA(info.VidmgrIpV4),
//	//MediaPort:  int(info.RtpPort),
//}
//
//_, err = media.SipPlayPush(stream)
//if err != nil {
//	fmt.Println("-----------------VidmgrGbsipChannelPlay  SipPlayPush----------------------")
//	return nil, errors.MediaSipPlayError.AddDetail("SipPush error:" + err.Error())
//}
//
////delete Stream
////streamRepo := db.NewVidmgrStreamRepo(l.ctx)
////streamFilter := db.VidmgrStreamFilter{
////	VidmgrID: pdev.VidmgrID,
////	Stream:   po.Stream,
////}
////streamRepo.DeleteByFilter(l.ctx, streamFilter)
////handle
//
//po.IsPlay = true
//po.Stream = stream.Stream
//if err := channelRepo.Update(l.ctx, po); err != nil {
//	er := errors.Fmt(err)
//	l.Errorf("%s  channel_id=%v  err=%v", utils.FuncName(), po.ChannelID, er)
//	return nil, er
//}

// 播放通道
func (l *SipChannelPlayLogic) SipChannelPlay(in *sip.SipChnPlayReq) (*sip.Response, error) {
	// todo: add your logic here and delete this line
	//media.Play(in.ChannelID, in.Replay, in.Start, in.End)
	return &sip.Response{}, nil
}
