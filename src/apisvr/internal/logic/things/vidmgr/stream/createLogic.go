package stream

import (
	"context"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.VidmgrStreamCreateReq) error {
	// todo: add your logic here and delete this line

	vidRep := &vid.VidmgrStreamCreateReq{
		RtpType: req.RtpType,
		StreamInfo: &vid.VidmgrStream{
			StreamName: req.StreamName,
			VidmgrID:   req.VidmgrID,

			App:    req.App,
			Stream: req.Stream,

			Identifier: req.Identifier,
			LocalIP:    req.LocalIP,
			LocalPort:  req.LocalPort,
			PeerIP:     req.PeerIP,
			PeerPort:   req.PeerPort,

			OriginStr:  req.OriginStr,
			OriginUrl:  req.OriginUrl,
			OriginType: req.OriginType,

			IsOnline:       req.IsOnline,
			IsRecordingHLS: req.IsRecordingHLS,
			IsRecordingMp4: req.IsRecordingMp4,
			IsShareChannel: req.IsShareChannel,
			IsAutoRecord:   req.IsAutoRecord,
			IsAutoPush:     req.IsAutoPush,
			IsPTZ:          req.IsPTZ,
			Desc:           utils.ToRpcNullString(req.Desc),
			Tags:           logic.ToTagsMap(req.Tags),
		},
	}
	if req.Vhost == "" {
		vidRep.StreamInfo.Vhost = clients.VHOSTNAME
	}
	_, err := l.svcCtx.VidmgrS.VidmgrStreamCreate(l.ctx, vidRep)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.VidmgrStreamCreate req=%v err=%v", utils.FuncName(), req, er)
		return er
	}

	return nil
}
