package areamanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreaInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.AreaInfoRepo
}

func NewAreaInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaInfoIndexLogic {
	return &AreaInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewAreaInfoRepo(ctx),
	}
}

// 获取区域信息列表
func (l *AreaInfoIndexLogic) AreaInfoIndex(in *sys.AreaInfoIndexReq) (*sys.AreaInfoIndexResp, error) {
	// todo: add your logic here and delete this line

	return &sys.AreaInfoIndexResp{}, nil
}
