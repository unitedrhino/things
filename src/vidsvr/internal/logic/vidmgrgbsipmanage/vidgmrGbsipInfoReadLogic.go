package vidmgrgbsipmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/vidsvr/internal/media"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidgmrGbsipInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidgmrGbsipInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidgmrGbsipInfoReadLogic {
	return &VidgmrGbsipInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取国标服务信息
func (l *VidgmrGbsipInfoReadLogic) VidgmrGbsipInfoRead(in *vid.VidmgrGbsipInfoReadReq) (*vid.VidmgrGbsipInfo, error) {
	// todo: add your logic here and delete this line

	pi, err := relationDB.NewVidmgrInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.VidmgrFilter{
		VidmgrIDs: []string{in.VidmgrID},
	})
	if err != nil {
		return nil, err
	}
	resp := &vid.VidmgrGbsipInfo{
		Region:       media.SipInfo.Region,
		CID:          media.SipInfo.CID,
		CNUM:         media.SipInfo.CNUM,
		DID:          media.SipInfo.DID,
		DNUM:         media.SipInfo.DNUM,
		LID:          media.SipInfo.LID,
		IP:           l.svcCtx.Config.Restconf.Host,
		Port:         l.svcCtx.Config.GbsipConf.UDP,
		MediaRtpPort: pi.RtpPort,
		MediaRtpIP:   pi.VidmgrIpV4,
	}
	return resp, nil
}
