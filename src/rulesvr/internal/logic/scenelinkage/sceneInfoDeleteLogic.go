package scenelinkagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type SceneInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSceneInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SceneInfoDeleteLogic {
	return &SceneInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SceneInfoDeleteLogic) SceneInfoDelete(in *rule.SceneInfoDeleteReq) (*rule.Response, error) {
	err := l.svcCtx.SceneRepo.Delete(l.ctx, in.Id)
	if err != nil { //如果是数据库错误
		return nil, errors.Database.AddDetail(err)
	}
	err = l.svcCtx.SceneTimerControl.Delete(in.Id)
	if err != nil {
		return nil, err
	}
	return &rule.Response{}, err
}
