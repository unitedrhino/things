package vidmgrstreammanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/clients"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/vidsvr/internal/common"
	"github.com/i-Things/things/service/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/vidsvr/internal/svc"
	"github.com/i-Things/things/service/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrStreamReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.VidmgrStreamRepo
}

func NewVidmgrStreamReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrStreamReadLogic {
	return &VidmgrStreamReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewVidmgrStreamRepo(ctx),
	}
}

// 获取流信息详情
func (l *VidmgrStreamReadLogic) VidmgrStreamRead(in *vid.VidmgrStreamReadReq) (*vid.VidmgrStream, error) {
	// todo: add your logic here and delete this line
	pi, err := l.PiDB.FindOneByFilter(l.ctx, relationDB.VidmgrStreamFilter{
		StreamIDs: []int64{in.StreamID},
	})
	if err != nil {
		return nil, err
	}
	vidInfo, err := relationDB.NewVidmgrInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.VidmgrFilter{
		VidmgrIDs: []string{pi.VidmgrID},
	})
	data := common.ToVidmgrStreamRpc(pi)
	//docker模式时，使用本地的RestCong的IP  端口还是使用配置的端口
	if vidInfo.MediasvrType == clients.MEDIA_DOCKER {
		data.MediaIP = l.svcCtx.Config.Restconf.Host
		data.MediaPort = l.svcCtx.Config.Mediakit.Port
	} else {
		data.MediaIP = utils.InetNtoA(vidInfo.VidmgrIpV4)
		data.MediaPort = vidInfo.VidmgrPort
	}
	return data, nil
}
