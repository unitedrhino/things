package scenelinkagelogic

import (
	"context"
	"github.com/i-Things/things/src/rulesvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"
	"github.com/zeromicro/go-zero/core/logx"
)

type SceneInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	SiDB *relationDB.SceneInfoRepo
}

func NewSceneInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SceneInfoReadLogic {
	return &SceneInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		SiDB:   relationDB.NewSceneInfoRepo(ctx),
	}
}

func (l *SceneInfoReadLogic) SceneInfoRead(in *rule.WithID) (*rule.SceneInfo, error) {
	pi, err := l.SiDB.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return ToScenePb(pi), nil
}
