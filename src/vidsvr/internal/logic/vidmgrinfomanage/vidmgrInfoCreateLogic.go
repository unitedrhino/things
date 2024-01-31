package vidmgrinfomanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/domain/deviceAuth"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/common"
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
		PiDB:   relationDB.NewVidmgrInfoRepo(ctx),
	}
}

// 新建服务
func (l *VidmgrInfoCreateLogic) VidmgrInfoCreate(in *vid.VidmgrInfo) (*vid.Response, error) {
	// todo: add your logic here and delete this line
	if in.VidmgrID == "" {
		randId := l.svcCtx.VidmgrID.GetSnowflakeId()
		in.VidmgrID = deviceAuth.GetStrProductID(randId)
	}
	//需要过滤相同配置 查找IP和端口相同的流服务，将错误
	filter := relationDB.VidmgrFilter{VidmgrIpV4: utils.InetAtoN(in.VidmgrIpV4), VidmgrPort: in.VidmgrPort}
	size, err := l.PiDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	if size > 0 {
		return nil, errors.MediaCreateError.AddDetailf("The MediaServer IP and Port is repeated:%d", size)
	}
	pi, err := common.ToVidmgrInfoDB(in)
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
