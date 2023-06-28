package scenelinkagelogic

import (
	"context"
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
	"github.com/i-Things/things/src/rulesvr/internal/logic"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type SceneInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSceneInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SceneInfoIndexLogic {
	return &SceneInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SceneInfoIndexLogic) SceneInfoIndex(in *rule.SceneInfoIndexReq) (*rule.SceneInfoIndexResp, error) {
	var (
		info []*rule.SceneInfo
		size int64
		err  error
	)
	filter := scene.InfoFilter{Name: in.Name, Status: in.Status,
		TriggerType: scene.TriggerType(in.TriggerType), AlarmID: in.AlarmID}
	size, err = l.svcCtx.SceneRepo.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	di, err := l.svcCtx.SceneRepo.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	info = make([]*rule.SceneInfo, 0, len(di))
	for _, v := range di {
		info = append(info, ToScenePb(v))
	}
	return &rule.SceneInfoIndexResp{List: info, Total: size}, nil
}
