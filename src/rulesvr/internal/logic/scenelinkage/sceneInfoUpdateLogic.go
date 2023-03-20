package scenelinkagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/rulesvr/internal/repo/mysql"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type SceneInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSceneInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SceneInfoUpdateLogic {
	return &SceneInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SceneInfoUpdateLogic) SceneInfoUpdate(in *rule.SceneInfo) (*rule.Response, error) {
	do, err := ToSceneDo(in)
	if err != nil {
		return nil, err
	}
	oldDo, err := l.svcCtx.SceneRepo.FindOneByName(l.ctx, do.Name)
	if err == nil && oldDo.ID != do.ID { //如果查到了并且和其他的场景重名了
		return nil, errors.Parameter.AddMsg("场景名字重复")
	}
	if err != mysql.ErrNotFound { //如果是数据库错误
		return nil, errors.Database.AddDetail(err)
	}
	err = do.Validate()
	if err != nil {
		return nil, err
	}
	if err = l.svcCtx.SceneRepo.Update(l.ctx, do); err != nil {
		return nil, err
	}
	if err = l.svcCtx.SceneTimerControl.Update(do); err != nil {
		return nil, err
	}

	return &rule.Response{}, err
}
