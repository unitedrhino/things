package vidmgrinfomanagelogic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/src/vidsvr/internal/common"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.VidmgrInfoRepo
}

func NewVidmgrInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrInfoReadLogic {
	return &VidmgrInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewVidmgrInfoRepo(ctx),
	}
}

// 获取服务信息详情
func (l *VidmgrInfoReadLogic) VidmgrInfoRead(in *vid.VidmgrInfoReadReq) (*vid.VidmgrInfo, error) {
	// todo: add your logic here and delete this line
	fmt.Printf("Vidsvr VidmgrInfoRead \n")
	pi, err := relationDB.NewVidmgrInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.VidmgrFilter{
		VidmgrIDs: []string{in.VidmgrID},
	})
	if err != nil {
		return nil, err
	}
	return common.ToVidmgrInfoRPC(pi), nil
}
