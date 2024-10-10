package rulelogic

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/things/service/udsvr/internal/domain/scene"
	"gitee.com/unitedrhino/things/service/udsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/udsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/udsvr/pb/ud"

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
	old, err := SceneInfoRead(l.ctx, l.svcCtx, in.Id, def.AuthReadWrite)
	if err != nil {
		return nil, err
	}
	if old.Tag == "deviceTiming" { //单设备定时
		uc := ctxs.GetUserCtx(l.ctx)
		di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
		})
		if err != nil {
			return nil, err
		}
		if uc.ProjectID != di.ProjectID {
			uc.ProjectID = di.ProjectID
			uc.IsAdmin = true
		}
	}
	var doUpdate bool
	if in.Name != "" {
		old.Name = in.Name
	}
	if in.Status != 0 {
		old.Status = in.Status
	}
	if in.IsCommon != 0 {
		old.IsCommon = in.IsCommon
	}
	if in.Desc != "" {
		old.Desc = in.Desc
	}
	if in.Logo != "" {
		old.Logo = in.Logo
	}
	if in.If != "" {
		old.UdSceneIf = newPo.UdSceneIf
		doUpdate = true
	}
	if in.When != "" {
		old.UdSceneWhen = newPo.UdSceneWhen
		doUpdate = true
	}
	if in.Then != "" {
		old.UdSceneThen = newPo.UdSceneThen
		doUpdate = true
	}
	if in.HeadImg != "" && in.IsUpdateHeadImg { //如果填了参数且不等于原来的,说明修改头像,需要处理
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessScene, oss.SceneHeadIng, fmt.Sprintf("%d/%s", old.ID, oss.GetFileNameWithPath(in.HeadImg)))
		path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(in.HeadImg, nwePath)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		}
		old.HeadImg = path
	}
	do := PoToSceneInfoDo(old)
	if doUpdate {
		err = do.Validate(NewSceneCheckRepo(l.ctx, l.svcCtx, do))
		if err != nil {
			return nil, err
		}
	}
	po := ToSceneInfoPo(do)
	po.SoftTime = old.SoftTime
	err = db.Update(l.ctx, po)
	if err != nil {
		return nil, err
	}
	if len(do.If.Triggers) == 0 && do.Type == scene.SceneTypeAuto && do.Status == def.True { //立即执行一次
		_, err = NewSceneManuallyTriggerLogic(l.ctx, l.svcCtx).SceneManuallyTrigger(&ud.WithID{Id: po.ID})
		if err != nil {
			return nil, err
		}
	}
	return &ud.Empty{}, err
}
