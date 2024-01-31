package vidmgrstreammanagelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/common"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrStreamIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.VidmgrStreamRepo
}

func NewVidmgrStreamIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrStreamIndexLogic {
	return &VidmgrStreamIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewVidmgrStreamRepo(ctx),
	}
}

// 获取流列表
func (l *VidmgrStreamIndexLogic) VidmgrStreamIndex(in *vid.VidmgrStreamIndexReq) (*vid.VidmgrStreamIndexResp, error) {
	// todo: add your logic here and delete this line
	var (
		info []*vid.VidmgrStream
		size int64
		err  error
	)
	filter := relationDB.VidmgrStreamFilter{
		VidmgrID:   in.VidmgrID,
		StreamIDs:  in.StreamIDs,
		App:        in.App,
		StreamName: in.StreamName,
		Stream:     in.Stream,
		Vhost:      in.Vhost,
		Identifier: in.Identifier,
		LocalIP:    utils.InetAtoN(in.LocalIP),
		LocalPort:  in.LocalPort,
		PeerIP:     utils.InetAtoN(in.PeerIP),
		PeerPort:   in.PeerPort,
	}
	size, err = l.PiDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	di, err := l.PiDB.FindByFilter(l.ctx, filter, common.ToPageInfoWithDefault(in.Page, &def.PageInfo{
		Page: 1, Size: 20,
		Orders: []def.OrderBy{{"created_time", def.OrderDesc}, {"stream_id", def.OrderDesc}},
	}))

	if err != nil {
		return nil, err
	}

	info = make([]*vid.VidmgrStream, 0, len(di))
	for _, v := range di {
		info = append(info, common.ToVidmgrStreamRpc(v))
	}
	fmt.Println("VidmgrStreamIndex:", info)
	return &vid.VidmgrStreamIndexResp{List: info, Total: size}, nil
}
