package scenelinkagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/events/topics"
	"github.com/i-Things/things/src/rulesvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type SceneInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	SiDB *relationDB.SceneInfoRepo
}

func NewSceneInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SceneInfoUpdateLogic {
	return &SceneInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		SiDB:   relationDB.NewSceneInfoRepo(ctx),
	}
}

func (l *SceneInfoUpdateLogic) SceneInfoUpdate(in *rule.SceneInfo) (*rule.Empty, error) {
	do, err := ToSceneDo(in)
	if err != nil {
		return nil, err
	}
	_, err = l.SiDB.FindOne(l.ctx, do.ID)
	if err != nil { //如果是数据库错误
		return nil, errors.Database.AddDetail(err)
	}
	err = do.Validate()
	if err != nil {
		return nil, err
	}
	if err = l.SiDB.Update(l.ctx, do); err != nil {
		return nil, err
	}
	if err = l.svcCtx.SceneTimerControl.Update(do); err != nil {
		return nil, err
	}
	if !l.svcCtx.SceneTimerControl.IsRunning() {
		l.svcCtx.Bus.Publish(l.ctx, topics.RuleSceneInfoUpdate, in.Id)
	}
	return &rule.Empty{}, err
}
