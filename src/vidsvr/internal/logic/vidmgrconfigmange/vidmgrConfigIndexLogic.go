package vidmgrconfigmangelogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/vidsvr/internal/logic"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrConfigIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.VidmgrConfigRepo
}

func NewVidmgrConfigIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrConfigIndexLogic {
	return &VidmgrConfigIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewVidmgrtConfigRepo(ctx),
	}
}

// 配置列表
func (l *VidmgrConfigIndexLogic) VidmgrConfigIndex(in *vid.VidmgrConfigIndexReq) (*vid.VidmgrConfigIndexResp, error) {
	// todo: add your logic here and delete this line
	//根据MediaserverID查找配置 返回所有的配置
	var (
		info []*vid.VidmgrConfig
		size int64
		err  error
	)
	filter := relationDB.VidmgrConfigFilter{MediaServerIds: in.MediaServerId}
	size, err = l.PiDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	di, err := l.PiDB.FindByFilter(l.ctx, filter, logic.ToPageInfoWithDefault(in.Page, &def.PageInfo{
		Page: 1, Size: 20,
		Orders: []def.OrderBy{{"created_time", def.OrderDesc}, {"id", def.OrderDesc}},
	}))
	if err != nil {
		return nil, err
	}
	info = make([]*vid.VidmgrConfig, 0, len(di))
	for _, v := range di {
		info = append(info, ToVidmgrConfig(v))
	}
	return &vid.VidmgrConfigIndexResp{
		List:  info,
		Total: size,
	}, nil
}
