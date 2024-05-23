package rulelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/oss"
	"github.com/i-Things/things/service/udsvr/internal/domain/scene"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

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
	if do.Status == 0 {
		do.Status = def.True
	}
	err := do.Validate(scene.ValidateRepo{Ctx: l.ctx, DeviceCache: l.svcCtx.DeviceCache, ProductSchemaCache: l.svcCtx.ProductSchemaCache})
	if err != nil {
		return nil, err
	}
	po := ToSceneInfoPo(do)
	err = relationDB.NewSceneInfoRepo(l.ctx).Insert(l.ctx, po)
	if err != nil {
		return nil, err
	}
	if in.HeadImg != "" && in.IsUpdateHeadImg { //如果填了参数且不等于原来的,说明修改头像,需要处理
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessScene, oss.SceneHeadIng, fmt.Sprintf("%d/%s", po.ID, oss.GetFileNameWithPath(in.HeadImg)))
		path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(in.HeadImg, nwePath)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		}
		po.HeadImg = path
		err = relationDB.NewSceneInfoRepo(l.ctx).UpdateHeadImg(l.ctx, po)
		if err != nil {
			l.Error(err)
		}
	}
	if len(do.If.Triggers) == 0 && do.Type == scene.SceneTypeAuto && do.Status == def.True { //立即执行一次
		_, err = NewSceneManuallyTriggerLogic(l.ctx, l.svcCtx).SceneManuallyTrigger(&ud.WithID{Id: po.ID})
		if err != nil {
			return nil, err
		}
	}
	return &ud.WithID{Id: po.ID}, nil
}
func GetSceneInfo(svcCtx svc.ServiceContext) func(ctx context.Context, sceneID int64) (info *scene.Info, err error) {
	return func(ctx context.Context, sceneID int64) (info *scene.Info, err error) {
		po, err := relationDB.NewSceneInfoRepo(ctx).FindOne(ctx, sceneID)
		if err != nil {
			return nil, err
		}
		return PoToSceneInfoDo(po), nil
	}
}
