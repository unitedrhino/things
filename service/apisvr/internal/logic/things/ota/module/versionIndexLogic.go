package module

import (
	"context"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VersionIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVersionIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VersionIndexLogic {
	return &VersionIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VersionIndexLogic) VersionIndex(req *types.OTAModuleIndexReq) (resp *types.OTAModuleVersionsIndexResp, err error) {
	// todo: add your logic here and delete this line

	return
}
