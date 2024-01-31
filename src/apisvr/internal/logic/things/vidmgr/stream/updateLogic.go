package stream

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.VidmgrStreamUpdateReq) error {
	// todo: add your logic here and delete this line
	vidReq := &vid.VidmgrStream{
		StreamID:       req.StreamID,
		VidmgrID:       req.VidmgrID,
		StreamName:     req.StreamName,
		App:            req.App,
		Stream:         req.Stream,
		Vhost:          req.Vhost,
		Identifier:     req.Identifier,
		LocalIP:        req.LocalIP,
		LocalPort:      req.LocalPort,
		PeerIP:         req.PeerIP,
		PeerPort:       req.PeerPort,
		OriginType:     req.OriginType,
		OriginUrl:      req.OriginUrl,
		OriginStr:      req.OriginStr,
		IsShareChannel: req.IsShareChannel,
		IsPTZ:          req.IsPTZ,
		IsAutoPush:     req.IsAutoPush,
		IsAutoRecord:   req.IsAutoRecord,
		IsRecordingMp4: req.IsRecordingMp4,
		IsRecordingHLS: req.IsRecordingHLS,
		IsOnline:       req.IsOnline,
		Desc:           utils.ToRpcNullString(req.Desc),
		Tags:           logic.ToTagsMap(req.Tags),
	}
	_, err := l.svcCtx.VidmgrS.VidmgrStreamUpdate(l.ctx, vidReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.VidmgrStreamUpdate req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
