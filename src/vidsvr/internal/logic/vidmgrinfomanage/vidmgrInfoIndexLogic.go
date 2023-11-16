package vidmgrinfomanagelogic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/vidsvr/internal/logic"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

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
		PiDB:   relationDB.NewVidmgrtInfoRepo(ctx),
	}
}

// 获取服务列表
func (l *VidmgrInfoIndexLogic) VidmgrInfoIndex(in *vid.VidmgrInfoIndexReq) (*vid.VidmgrInfoIndexResp, error) {
	// todo: add your logic here and delete this line
	fmt.Printf("Vidsvr VidmgrInfoIndex \n")

	var (
		info []*vid.VidmgrInfo
		size int64
		err  error
		piDB = relationDB.NewVidmgrtInfoRepo(l.ctx)
	)
	filter := relationDB.VidmgrFilter{VidmgrType: in.VidmgrType, VidmgrName: in.VidmgrtName, Tags: in.Tags, VidmgrIDs: in.VidmgrIDs}
	size, err = piDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	di, err := piDB.FindByFilter(l.ctx, filter, logic.ToPageInfoWithDefault(in.Page, &def.PageInfo{
		Page: 1, Size: 20,
		Orders: []def.OrderBy{{"created_time", def.OrderDesc}, {"id", def.OrderDesc}},
	}))

	if err != nil {
		return nil, err
	}

	info = make([]*vid.VidmgrInfo, 0, len(di))
	for _, v := range di {
		info = append(info, ToVidmgrInfo(l.ctx, v, l.svcCtx))
	}

	return &vid.VidmgrInfoIndexResp{List: info, Total: size}, nil
}
