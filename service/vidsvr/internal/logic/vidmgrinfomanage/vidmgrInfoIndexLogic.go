package vidmgrinfomanagelogic

import (
	"context"
	"gitee.com/i-Things/share/def"
	"github.com/i-Things/things/service/vidsvr/internal/common"
	"github.com/i-Things/things/service/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/vidsvr/internal/svc"
	"github.com/i-Things/things/service/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.VidmgrInfoRepo
}

func NewVidmgrInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrInfoIndexLogic {
	return &VidmgrInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewVidmgrInfoRepo(ctx),
	}
}

// 获取服务列表
func (l *VidmgrInfoIndexLogic) VidmgrInfoIndex(in *vid.VidmgrInfoIndexReq) (*vid.VidmgrInfoIndexResp, error) {
	// todo: add your logic here and delete this line
	var (
		info []*vid.VidmgrInfo
		size int64
		err  error
	)
	filter := relationDB.VidmgrFilter{VidmgrIDs: in.VidmgrIDs, VidmgrName: in.VidmgrtName, VidmgrType: in.VidmgrType, Tags: in.Tags}
	size, err = l.PiDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	di, err := l.PiDB.FindByFilter(l.ctx, filter, common.ToPageInfoWithDefault(in.Page, &def.PageInfo{
		Page: 1, Size: 20,
		Orders: []def.OrderBy{{"created_time", def.OrderDesc}, {"vidmgr_id", def.OrderDesc}},
	}))

	if err != nil {
		return nil, err
	}

	info = make([]*vid.VidmgrInfo, 0, len(di))
	for _, v := range di {
		info = append(info, common.ToVidmgrInfoRPC(v))
	}

	return &vid.VidmgrInfoIndexResp{List: info, Total: size}, nil
}
