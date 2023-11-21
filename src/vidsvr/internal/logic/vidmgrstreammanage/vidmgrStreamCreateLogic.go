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

// 流添加
func (l *VidmgrStreamCreateLogic) VidmgrStreamCreate(in *vid.VidmgrStream) (*vid.Response, error) {
	// todo: add your logic here and delete this line

	//插入数据之前确定流服务是否存在
	pi, err := relationDB.NewVidmgrInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.VidmgrFilter{
		VidmgrIDs: []string{in.VidmgrID},
	})
	if err != nil {
		return nil, err
	}
	//插入数据
	if pi.VidmgrID != "" {
		err = l.PiDB.Insert(l.ctx, ToDbConvVidmgrStream(in))
		if err != nil {
			l.Errorf("%s.Insert err=%+v", utils.FuncName(), err)
			return nil, err
		}
	}

	return &vid.Response{}, nil
}
