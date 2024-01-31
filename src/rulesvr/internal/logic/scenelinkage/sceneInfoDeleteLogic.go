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

type SceneInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	SiDB *relationDB.SceneInfoRepo
}

func NewSceneInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SceneInfoDeleteLogic {
	return &SceneInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		SiDB:   relationDB.NewSceneInfoRepo(ctx),
	}
}

func (l *SceneInfoDeleteLogic) SceneInfoDelete(in *rule.WithID) (*rule.Empty, error) {
	err := l.SiDB.Delete(l.ctx, in.Id)
	if err != nil { //如果是数据库错误
		return nil, errors.Database.AddDetail(err)
	}
	err = l.svcCtx.SceneTimerControl.Delete(in.Id)
	if err != nil {
		return nil, err
	}
	if !l.svcCtx.SceneTimerControl.IsRunning() {
		l.svcCtx.Bus.Publish(l.ctx, topics.RuleSceneInfoDelete, in.Id)
	}
	return &rule.Empty{}, err
}
