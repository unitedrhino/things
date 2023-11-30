package stream

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/vidsvr/pb/vid"
)

func VidmgrStreamToApi(v *vid.VidmgrStream) *types.VidmgrStream {
	return &types.VidmgrStream{
		StreamID:         v.StreamID,
		Protocol:         v.Protocol,
		ReaderCount:      v.ReaderCount,
		TotalReaderCount: v.TotalReaderCount,

		StreamCommon: types.StreamCommon{
			App:            v.App,
			Stream:         v.Stream,
			StreamName:     v.StreamName,
			VidmgrID:       v.VidmgrID,
			Vhost:          v.Vhost,
			Identifier:     v.Identifier,
			LocalIP:        v.LocalIP,
			LocalPort:      v.LocalPort,
			PeerIP:         v.PeerIP,
			PeerPort:       v.PeerPort,
			OriginType:     v.OriginType,
			OriginUrl:      v.OriginUrl,
			OriginStr:      v.OriginStr,
			IsShareChannel: v.IsShareChannel,
			IsPTZ:          v.IsPTZ,
			IsAutoPush:     v.IsAutoPush,
			IsAutoRecord:   v.IsAutoRecord,
			IsRecordingMp4: v.IsRecordingMp4,
			IsRecordingHLS: v.IsRecordingHLS,
			IsOnline:       v.IsOnline,
			Tags:           logic.ToTagsType(v.Tags),
			Desc:           utils.ToNullString(v.Desc),
		},
		Tracks: TovidmgrTracksApi(v.Tracks),
	}
}

func ToVidmgrTrackApi(in *vid.StreamTrack) *types.StreamTrack {
	pi := &types.StreamTrack{
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

func TovidmgrTracksApi(pi []*vid.StreamTrack) []types.StreamTrack {
	if len(pi) > 0 {
		info := make([]types.StreamTrack, 0, len(pi))
		for _, v := range pi {
			info = append(info, *ToVidmgrTrackApi(v))
		}
		return info
	}
	return nil
}
