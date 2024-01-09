package rulelogic

import (
	"context"
	"github.com/i-Things/things/src/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/udsvr/internal/svc"
	"github.com/i-Things/things/src/udsvr/pb/ud"

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

func (l *SceneInfoUpdateLogic) SceneInfoUpdate(in *ud.SceneInfo) (*ud.Empty, error) {
	newPo := ToSceneInfoPo(ToSceneInfoDo(in))
	db := relationDB.NewSceneInfoRepo(l.ctx)
	old, err := db.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.Name != "" {
		old.Name = in.Name
	}
	if in.Status != 0 {
		old.Status = in.Status
	}
	if in.Desc != "" {
		old.Desc = in.Desc
	}
	if in.Trigger != "" {
		old.UdSceneTrigger = newPo.UdSceneTrigger
	}
	if in.When != "" {
		old.UdSceneWhen = newPo.UdSceneWhen
	}
	if in.Then != "" {
		old.UdSceneThen = newPo.UdSceneThen
	}
	err = PoToSceneInfoDo(old).Validate()
	if err != nil {
		return nil, err
	}
	err = db.Update(l.ctx, old)
	return &ud.Empty{}, err
}
