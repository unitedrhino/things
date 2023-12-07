package vidmgrstreammanagelogic

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/pb/vid"
	"time"
)

func ToDbConvVidmgrStream(in *vid.VidmgrStream) *relationDB.VidmgrStream {
	info := make([]*relationDB.StreamTrack, 0, len(in.Tracks))
	for _, v := range in.Tracks {
		info = append(info, ToDbVidmgrStreamTrack(v))
	}
	pi := &relationDB.VidmgrStream{
		VidmgrID:   in.VidmgrID,
		StreamName: in.StreamName,

		App: in.App,
		//Schema: in.Schema,
		Protocol: in.Protocol,
		Stream:   in.Stream,
		Vhost:    in.Vhost,

		Identifier: in.Identifier,
		LocalIP:    utils.InetAtoN(in.LocalIP),
		LocalPort:  in.LocalPort,
		PeerIP:     utils.InetAtoN(in.PeerIP),
		PeerPort:   in.PeerPort,

		OriginType:       in.OriginType,
		OriginStr:        in.OriginStr,
		OriginUrl:        in.OriginUrl,
		ReaderCount:      in.ReaderCount,
		TotalReaderCount: in.TotalReaderCount,
		Tracks:           info,
		IsRecordingHLS:   in.IsRecordingHLS,
		IsRecordingMp4:   in.IsRecordingMp4,
		IsShareChannel:   in.IsShareChannel,
		IsAutoPush:       in.IsAutoPush,
		IsAutoRecord:     in.IsAutoRecord,
		IsPTZ:            in.IsPTZ,
		IsOnline:         in.IsOnline,
		//LastLogin:        time.Unix(in.LastLogin, 0),
		Desc: in.Desc.GetValue(),
	}
	if in.Tags == nil {
		in.Tags = map[string]string{}
	}
	pi.Tags = in.Tags

	return pi
}

func ToDbVidmgrStreamTrack(in *vid.StreamTrack) *relationDB.StreamTrack {
	pi := &relationDB.StreamTrack{
		Channels:    in.Channels,
		CodecId:     in.CodecId,
		CodecIdName: in.CodecIdName,
		CodecType:   in.CodecType,
		Ready:       in.Ready,
		SampleBit:   in.SampleBit,
		SampleRate:  in.SampleRate,
		Fps:         in.Fps,
		Height:      in.Height,
		Width:       in.Width,
	}
	return pi
}

func ToRpcVidmgrStreamTrack(in *relationDB.StreamTrack) *vid.StreamTrack {
	pi := &vid.StreamTrack{
		Channels:    in.Channels,
		CodecId:     in.CodecId,
		CodecIdName: in.CodecIdName,
		CodecType:   in.CodecType,
		Ready:       in.Ready,
		SampleBit:   in.SampleBit,
		SampleRate:  in.SampleRate,
		Fps:         in.Fps,
		Height:      in.Height,
		Width:       in.Width,
	}
	return pi
}

func ToRpcConvVidmgrStream(in *relationDB.VidmgrStream) *vid.VidmgrStream {

	info := make([]*vid.StreamTrack, 0, len(in.Tracks))
	for _, v := range in.Tracks {
		info = append(info, ToRpcVidmgrStreamTrack(v))
	}
	pi := &vid.VidmgrStream{
		StreamID:   in.StreamID,
		VidmgrID:   in.VidmgrID,
		StreamName: in.StreamName,

		App:      in.App,
		Protocol: in.Protocol,
		Stream:   in.Stream,
		Vhost:    in.Vhost,

		Identifier: in.Identifier,
		LocalIP:    utils.InetNtoA(in.LocalIP),
		LocalPort:  in.LocalPort,
		PeerIP:     utils.InetNtoA(in.PeerIP),
		PeerPort:   in.PeerPort,

		OriginType:       in.OriginType,
		OriginStr:        in.OriginStr,
		OriginUrl:        in.OriginUrl,
		ReaderCount:      in.ReaderCount,
		TotalReaderCount: in.TotalReaderCount,
		Tracks:           info,
		IsRecordingHLS:   in.IsRecordingHLS,
		IsRecordingMp4:   in.IsRecordingMp4,
		IsShareChannel:   in.IsShareChannel,
		IsAutoPush:       in.IsAutoPush,
		IsAutoRecord:     in.IsAutoRecord,
		IsPTZ:            in.IsPTZ,
		IsOnline:         in.IsOnline,
		Desc:             utils.ToRpcNullString(&in.Desc),
		Tags:             in.Tags,
	}

	return pi
}

func setPoByPb(old *relationDB.VidmgrStream, data *vid.VidmgrStream) error {
	if data.StreamName != "" {
		old.StreamName = data.StreamName
	}
	if data.App != "" {
		old.App = data.App
	}
	if data.Protocol != 0 {
		old.Protocol = data.Protocol
	}
	if data.Stream != "" {
		old.Stream = data.Stream
	}
	if data.Vhost != "" {
		old.Vhost = data.Vhost
	}
	if data.Identifier != "" {
		old.Identifier = data.Identifier
	}
	if data.LocalIP != "" {
		old.LocalIP = utils.InetAtoN(data.LocalIP)
	}
	if data.LocalPort != 0 {
		old.LocalPort = data.LocalPort
	}
	if data.PeerIP != "" {
		old.PeerIP = utils.InetAtoN(data.PeerIP)
	}
	if data.PeerPort != 0 {
		old.PeerPort = data.PeerPort
	}
	if data.OriginType != 0 {
		old.OriginType = data.OriginType
	}
	if data.OriginStr != "" {
		old.OriginStr = data.OriginStr
	}
	if data.OriginUrl != "" {
		old.OriginUrl = data.OriginUrl
	}
	if data.VidmgrID != "" {
		old.VidmgrID = data.VidmgrID
	}
	if len(data.Tracks) > 0 {
		info := make([]*relationDB.StreamTrack, 0, len(data.Tracks))
		for _, v := range data.Tracks {
			info = append(info, ToDbVidmgrStreamTrack(v))
		}
		old.Tracks = info
	}
	if data.LastLogin != 0 {
		old.LastLogin.Valid = true
		old.LastLogin.Time = time.Unix(data.LastLogin, 0)
	}

	old.IsRecordingMp4 = data.IsRecordingMp4
	old.IsRecordingHLS = data.IsRecordingHLS
	old.IsShareChannel = data.IsShareChannel
	old.IsAutoPush = data.IsAutoPush
	old.IsAutoRecord = data.IsAutoRecord
	old.IsPTZ = data.IsPTZ
	old.IsOnline = old.IsOnline
	old.ReaderCount = data.ReaderCount
	old.TotalReaderCount = data.TotalReaderCount
	return nil
}
