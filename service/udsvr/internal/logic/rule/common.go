package rulelogic

import (
	"context"
	"github.com/i-Things/things/service/udsvr/internal/domain/scene"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/udsvr/internal/svc"
)

func NewSceneValidateRepo(ctx context.Context, svcCtx *svc.ServiceContext) scene.ValidateRepo {
	return scene.ValidateRepo{
		Ctx: ctx, DeviceCache: svcCtx.DeviceCache,
		ProductCache:       svcCtx.ProductCache,
		ProductSchemaCache: svcCtx.ProductSchemaCache,
		ProjectCache:       svcCtx.ProjectCache,
		UserShareCache:     svcCtx.UserShareCache,
		GetSceneInfo:       GetSceneInfo,
	}
}
func GetSceneInfo(ctx context.Context, sceneID int64) (info *scene.Info, err error) {
	po, err := relationDB.NewSceneInfoRepo(ctx).FindOne(ctx, sceneID)
	if err != nil {
		return nil, err
	}
	return PoToSceneInfoDo(po), nil
}
