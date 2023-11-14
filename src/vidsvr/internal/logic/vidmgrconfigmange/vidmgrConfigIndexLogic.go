package vidmgrconfigmangelogic

import (
	"context"
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
func (l *VidmgrConfigIndexLogic) VidmgrConfigIndex(in *vid.VidmgrConfigReadReq) (*vid.VidmgrConfigIndexResp, error) {
	// todo: add your logic here and delete this line

	return &vid.VidmgrConfigIndexResp{}, nil
}
