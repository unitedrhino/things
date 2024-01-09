package rulelogic

import (
	"context"
	"github.com/i-Things/things/src/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/udsvr/internal/svc"
	"github.com/i-Things/things/src/udsvr/pb/ud"

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

// 场景
func (l *SceneInfoCreateLogic) SceneInfoCreate(in *ud.SceneInfo) (*ud.WithID, error) {
	do := ToSceneInfoDo(in)
	err := do.Validate()
	if err != nil {
		return nil, err
	}
	//校验数据权限 todo

	po := ToSceneInfoPo(do)
	relationDB.NewSceneInfoRepo(l.ctx).Insert(l.ctx, po)

	return &ud.WithID{Id: po.ID}, nil
}
