package vidmgrconfigmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrConfigDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.VidmgrConfigRepo
}

func NewVidmgrConfigDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrConfigDeleteLogic {
	return &VidmgrConfigDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewVidmgrConfigRepo(ctx),
	}
}

// 删除配置
func (l *VidmgrConfigDeleteLogic) VidmgrConfigDelete(in *vid.VidmgrConfigDeleteReq) (*vid.Response, error) {
	// todo: add your logic here and delete this line
	err := l.PiDB.DeleteByFilter(l.ctx, relationDB.VidmgrConfigFilter{VidmgrIDs: []string{in.GeneralMediaServerId}})
	if err != nil {
		l.Errorf("%s.Delete err=%v", utils.FuncName(), utils.Fmt(err))
		return nil, err
	}
	return &vid.Response{}, nil
}
