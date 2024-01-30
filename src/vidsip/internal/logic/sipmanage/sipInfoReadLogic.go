package sipmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/vidsip/internal/media"

	"github.com/i-Things/things/src/vidsip/internal/svc"
	"github.com/i-Things/things/src/vidsip/pb/sip"

	"github.com/zeromicro/go-zero/core/logx"
)

type SipInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSipInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SipInfoReadLogic {
	return &SipInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取国标服务信息
func (l *SipInfoReadLogic) SipInfoRead(in *sip.SipInfoReadReq) (*sip.SipInfo, error) {
	// todo: add your logic here and delete this line
	resp := &sip.SipInfo{
		Region: media.SipInfo.Region,
		CID:    media.SipInfo.CID,
		CNUM:   media.SipInfo.CNUM,
		DID:    media.SipInfo.DID,
		DNUM:   media.SipInfo.DNUM,
		LID:    media.SipInfo.LID,
		IP:     media.SipInfo.SipIp,
		Port:   int64(media.SipInfo.SipPort),
	}
	return resp, nil
}
