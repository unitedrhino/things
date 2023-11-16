package vidmgrinfomanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/domain/deviceAuth"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.VidmgrInfoRepo
}

func NewVidmgrInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrInfoCreateLogic {
	return &VidmgrInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewVidmgrtInfoRepo(ctx),
	}
}

// 新建服务
func (l *VidmgrInfoCreateLogic) VidmgrInfoCreate(in *vid.VidmgrInfo) (*vid.Response, error) {
	// todo: add your logic here and delete this line
	if in.VidmgrID == "" {
		randId := l.svcCtx.VidmgrID.GetSnowflakeId()
		in.VidmgrID = deviceAuth.GetStrProductID(randId)
	}
	pi, err := ConvVidmgrPbToPo(in)
	if err != nil {
		return nil, err
	}
	err = l.PiDB.Insert(l.ctx, pi)
	if err != nil {
		l.Errorf("%s.Insert err=%+v", utils.FuncName(), err)
		return nil, err
	}
	return &vid.Response{}, nil
}
