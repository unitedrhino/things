package vidmgrconfigmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/vidsvr/internal/common"
	"github.com/i-Things/things/service/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/vidsvr/internal/svc"
	"github.com/i-Things/things/service/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrConfigCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.VidmgrConfigRepo
}

func NewVidmgrConfigCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrConfigCreateLogic {
	return &VidmgrConfigCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewVidmgrConfigRepo(ctx),
	}
}

// 新建配置
func (l *VidmgrConfigCreateLogic) VidmgrConfigCreate(in *vid.VidmgrConfig) (*vid.Response, error) {
	// todo: add your logic here and delete this line
	dbConfig := common.ToVidmgrConfigDB(in)
	err := l.PiDB.Insert(l.ctx, dbConfig)
	if err != nil {
		l.Errorf("%s.Insert err=%+v", utils.FuncName(), err)
		return nil, err
	}
	return &vid.Response{}, nil
}
