package vidmgrmangelogic

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

// 服务管理
func (l *VidmgrInfoCreateLogic) VidmgrInfoCreate(in *vid.VidmgrInfo) (*vid.Response, error) {
	// todo: add your logic here and delete this line
	if in.VidmgrID == "" {
		randId := l.svcCtx.VidmgrID.GetSnowflakeId()
		in.VidmgrID = deviceAuth.GetStrProductID(randId)
	}
	pi, err := l.ConvVidmgrPbToPo(in)
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

/*
根据用户的输入生成对应的数据库数据
*/
func (l *VidmgrInfoCreateLogic) ConvVidmgrPbToPo(in *vid.VidmgrInfo) (*relationDB.VidmgrInfo, error) {
	pi := &relationDB.VidmgrInfo{
		VidmgrID:     in.VidmgrID,
		VidmgrName:   in.VidmgrName,
		VidmgrIpV4:   utils.InetAtoN(in.VidmgrIpV4),
		VidmgrPort:   in.VidmgrPort,
		VidmgrStatus: in.VidmgrStatus,
		VidmgrSecret: in.VidmgrSecret,
		VidmgrType:   in.VidmgrType,
		Desc:         in.Desc.GetValue(),
	}
	if in.Tags == nil {
		in.Tags = map[string]string{}
	}
	pi.Tags = in.Tags
	return pi, nil
}
