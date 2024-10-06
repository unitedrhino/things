package rulelogic

import (
	"context"
	"gitee.com/i-Things/things/service/udsvr/internal/domain/scene"
	"gitee.com/i-Things/things/service/udsvr/internal/repo/relationDB"
	"gitee.com/i-Things/things/service/udsvr/internal/svc"
)

func NewSceneCheckRepo(ctx context.Context, svcCtx *svc.ServiceContext, do *scene.Info) scene.CheckRepo {
	return scene.CheckRepo{
		Ctx: ctx, DeviceCache: svcCtx.DeviceCache,
		ProductCache:       svcCtx.ProductCache,
		ProductSchemaCache: svcCtx.ProductSchemaCache,
		ProjectCache:       svcCtx.ProjectCache,
		UserShareCache:     svcCtx.UserShareCache,
		DeviceMsg:          svcCtx.DeviceMsg,
		Common:             svcCtx.SysCommon,
		Info:               do,
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
