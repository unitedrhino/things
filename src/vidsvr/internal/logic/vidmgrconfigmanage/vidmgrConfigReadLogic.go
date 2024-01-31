package vidmgrconfigmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/common"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrConfigReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.VidmgrConfigRepo
}

func NewVidmgrConfigReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrConfigReadLogic {
	return &VidmgrConfigReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewVidmgrConfigRepo(ctx),
	}
}

// 获取配置信息详情
func (l *VidmgrConfigReadLogic) VidmgrConfigRead(in *vid.VidmgrConfigReadReq) (*vid.VidmgrConfig, error) {
	// todo: add your logic here and delete this line
	pi, err := relationDB.NewVidmgrConfigRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.VidmgrConfigFilter{
		VidmgrIDs: []string{in.MediaServerId},
	})
	if err != nil {
		l.Errorf("%s.Delete err=%v", utils.FuncName(), utils.Fmt(err))
		return nil, err
	}
	return common.ToVidmgrConfigRpc(pi), nil
}
