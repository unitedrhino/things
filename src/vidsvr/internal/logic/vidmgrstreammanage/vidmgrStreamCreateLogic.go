package vidmgrstreammanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrStreamCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.VidmgrStreamRepo
}

func NewVidmgrStreamCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrStreamCreateLogic {
	return &VidmgrStreamCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewVidmgrStreamRepo(ctx),
	}
}

// 流添加 拉流添加接口
func (l *VidmgrStreamCreateLogic) VidmgrStreamCreate(in *vid.VidmgrStream) (*vid.Response, error) {
	// todo: add your logic here and delete this line

	err := l.PiDB.Insert(l.ctx, ToDbConvVidmgrStream(in))
	if err != nil {
		l.Errorf("%s.Insert err=%+v", utils.FuncName(), err)
		return nil, err
	}

	return &vid.Response{}, nil
}
