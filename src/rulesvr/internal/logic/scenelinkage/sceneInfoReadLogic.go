package scenelinkagelogic

import (
	"context"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type SceneInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSceneInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SceneInfoReadLogic {
	return &SceneInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SceneInfoReadLogic) SceneInfoRead(in *rule.SceneInfoReadReq) (*rule.SceneInfo, error) {
	// todo: add your logic here and delete this line

	return &rule.SceneInfo{}, nil
}
