package scenelinkagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/rulesvr/internal/repo/mysql"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type SceneInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSceneInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SceneInfoCreateLogic {
	return &SceneInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SceneInfoCreateLogic) SceneInfoCreate(in *rule.SceneInfo) (*rule.Response, error) {
	do, err := ToSceneDo(in)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.SceneRepo.FindOneByName(l.ctx, do.Name)
	if err == nil {
		return nil, errors.Parameter.AddMsg("场景名字重复")
	}
	if err != mysql.ErrNotFound {
		return nil, errors.Database.AddDetail(err)
	}
	err = do.Validate()
	if err != nil {
		return nil, err
	}
	id, err := l.svcCtx.SceneRepo.Insert(l.ctx, do)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.SceneTimerControl.Create(do)
	if err != nil {
		return nil, err
	}
	return &rule.Response{Id: id}, err
}
