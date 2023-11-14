package info

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActiveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewActiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActiveLogic {
	return &ActiveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ActiveLogic) Active(req *types.VidmgrInfoActiveReq) error {
	// todo: add your logic here and delete this line
	//read VidmgrInfo table and update table
	vidTmp, err := l.svcCtx.VidmgrM.VidmgrInfoRead(l.ctx, &vid.VidmgrInfoReadReq{
		VidmgrtID: req.VidmgrID,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ManageVidmgr req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	if vidTmp.VidmgrStatus == def.DeviceStatusInactive {
		//UPDATE
		vidReq := &vid.VidmgrInfo{
			VidmgrName:   vidTmp.VidmgrName,
			VidmgrID:     vidTmp.VidmgrID,
			VidmgrIpV4:   vidTmp.VidmgrIpV4,
			VidmgrPort:   vidTmp.VidmgrPort,
			VidmgrType:   vidTmp.VidmgrType,
			VidmgrSecret: vidTmp.VidmgrSecret,
			VidmgrStatus: def.DeviceStatusOffline,
		}
		vidTmp.VidmgrStatus = def.DeviceStatusOffline

		_, err := l.svcCtx.VidmgrM.VidmgrInfoUpdate(l.ctx, vidReq)
		if err != nil {
			er := errors.Fmt(err)
			l.Errorf("%s.rpc.ManageVidmgr req=%v err=%v", utils.FuncName(), req, er)
			return er
		}
		//更新之后需要配置流媒体服务
		//set default
	}

	return nil
}
